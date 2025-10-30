package notifier

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

const (
	maxRetries = 5
	baseDelay  = 30 * time.Second
)

func getRetryCount(headers amqp091.Table) int {
	raw := headers["x-retry-count"]
	val, ok := raw.(int)
	// p.s. в случае если появится несколько сервисов,
	// реализовать type-switch и учитывать также например string (для большей безопасности, что мы точно не потеряем этот x-retry-count и не пойдем retry-ить с начала с 0)
	if !ok {
		return 0
	}
	return val
}

func (n *Notifier) updateDB(ctx context.Context, notificationID int64, status models.Status) error {
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Millisecond * 500,
		Backoff:  2,
	}
	updateFn := func() error { return n.usecase.UpdateNotificationStatus(ctx, notificationID, status) }
	err := retry.Do(updateFn, strategy)
	if err != nil {
		return fmt.Errorf("update notification status: %w", err)
	}
	return nil
}

func (n *Notifier) calculateExponentialDelay(retryCount int) time.Duration {
	// 30s → 60s → 120s → 240s → 480s
	return time.Duration(math.Pow(2, float64(retryCount))) * baseDelay
}

func (n *Notifier) processRetry(ctx context.Context, delivery amqp091.Delivery, notificationID int64, channelType models.Channel) {
	retryCount := getRetryCount(delivery.Headers)
	if retryCount == maxRetries {
		err := n.updateDB(ctx, notificationID, models.StatusFailed)
		if err != nil {
			zlog.Logger.Error().
				Err(err).
				Str("path", "processConsume n.updateDB").
				Int64("notification_id", notificationID)
			// мб кинуть алерт, тк состояние бд неконсистентно
		}
		zlog.Logger.Warn().Msgf("notification %d failed after %d retries", notificationID, retryCount)
		delivery.Ack(false)
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
		delivery.Nack(false, true) // пробуем еще раз
		return
	}

	delivery.Ack(false)

}
