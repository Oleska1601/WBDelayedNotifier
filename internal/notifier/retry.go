package notifier

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

const (
	maxRetries = 5
	baseDelay  = 1 * time.Second
)

func getRetryCount(headers amqp091.Table) int {
	raw, exists := headers["x-retry-count"]
	if !exists {
		return 0
	}

	str := fmt.Sprintf("%v", raw)
	count, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return count
}

func (n *Notifier) updateDB(ctx context.Context, updateNotification models.UpdateNotification) error {
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Millisecond * 500,
		Backoff:  2,
	}
	// ф-ия с retry, поскольку необходимо обеспечить консистентность между статусом в БД и фактической отправкой:
	// - после исчерпания лимита попыток обновляет статус на 'failed'
	// - при успешной отправке обновляет статус на 'sent' и устанавливает время отправки
	updateFn := func() error { return n.usecase.UpdateNotification(ctx, updateNotification) }
	err := retry.Do(updateFn, strategy)
	if err != nil {
		return fmt.Errorf("n.usecase.UpdateNotificationStatus: %w", err)
	}
	return nil
}

func (n *Notifier) calculateExponentialDelay(retryCount int) time.Duration {
	return time.Duration(math.Pow(2, float64(retryCount))) * baseDelay
}

func (n *Notifier) processRetry(ctx context.Context, delivery amqp091.Delivery, notificationID int64, channelType models.Channel) {
	retryCount := getRetryCount(delivery.Headers)
	if retryCount == maxRetries {
		updateNotification := models.UpdateNotification{
			ID:     notificationID,
			Status: models.StatusFailed,
		}
		err := n.updateDB(ctx, updateNotification)
		if err != nil {
			zlog.Logger.Error().
				Err(err).
				Str("path", "processConsume n.updateDB").
				Int64("notification_id", notificationID)
		}

		zlog.Logger.Error().Msgf("notification %d failed after %d retries", notificationID, retryCount)
		n.safeAck(false, delivery)
		return
	}

	delay := n.calculateExponentialDelay(retryCount)

	var routingKey string
	if channelType == models.ChannelEmail {
		routingKey = n.cfg.RetryQueueEmail.RoutingKey
	} else {
		routingKey = n.cfg.RetryQueueTg.RoutingKey
	}

	err := n.channel.Publish(
		n.cfg.RetryExchange.Name, routingKey, false, false,
		amqp091.Publishing{
			Body: delivery.Body,
			Headers: amqp091.Table{
				"x-retry-count": retryCount + 1,
				"x-delay":       delay,
			},
		},
	)

	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "processRetry n.channel.Publish").
			Int64("notification_id", notificationID)
		n.safeNack(false, true, delivery) // пробуем еще раз
		return
	}

	n.safeAck(false, delivery)

}
