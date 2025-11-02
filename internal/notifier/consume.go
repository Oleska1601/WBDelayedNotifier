package notifier

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/zlog"
)

const (
	workers      = 5
	retryWorkers = 2
)

func (n *Notifier) consume(ctx context.Context, workers int, queueName string, channelType models.Channel) {
	msgChan, err := n.channel.Consume(
		queueName,
		"", false, false, false, false, nil,
	)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "Consume n.channel.Consume").
			Msgf("failed to consume from %s queue", queueName)
		return
	}

	wg := &sync.WaitGroup{}
	for range workers {
		wg.Add(1)
		go n.worker(ctx, msgChan, wg, channelType)
	}
	wg.Wait()
	zlog.Logger.Info().
		Str("queue", queueName).
		Msg("all workers stopped")

}

func (n *Notifier) worker(ctx context.Context, msgChan <-chan amqp091.Delivery, wg *sync.WaitGroup, channelType models.Channel) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-msgChan:
			if !ok {
				return
			}
			n.processConsume(ctx, delivery, channelType)
		}
	}
}

func (n *Notifier) processConsume(ctx context.Context, delivery amqp091.Delivery, channelType models.Channel) {
	sender, exists := n.senders[channelType]
	if !exists {
		zlog.Logger.Error().
			Type("channel_type", channelType).
			Msgf("channel type %s is not registered", channelType)

		n.safeNack(false, false, delivery)
		return
	}

	var notification models.NotificationMessage
	// 1
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "process json.Unmarshal").
			Msg("failed to unmarshal json")

		n.safeNack(false, false, delivery)
		return
	}

	// 2
	currentStatus, err := n.usecase.GetNotificationStatus(ctx, notification.ID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "processConsume n.usecase.GetNotificationStatus").
			Msg("failed to get notification status")

		n.safeNack(false, false, delivery)
		return
	}

	// 3
	if currentStatus == models.StatusCancelled {
		zlog.Logger.Info().
			Str("path", "processConsume").
			Int64("notification_id", notification.ID).
			Msg("send notification is cancelled")
		n.safeAck(false, delivery) // сообщение отменено - доставлять никуда не надо -> его обработка завершена
		return
	}

	// 4
	if err := sender.Send(notification.Recipient, notification.Message); err != nil {
		n.processRetry(ctx, delivery, notification.ID, channelType)
		return
	}
	now := time.Now()
	updateNotification := models.UpdateNotification{
		ID:     notification.ID,
		SentAt: &now,
		Status: models.StatusSent,
	}

	err = n.updateDB(ctx, updateNotification)

	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "processConsume n.updateDB").
			Int64("notification_id", notification.ID)
	}
	zlog.Logger.Info().Msgf("notification %d is processed successful", notification.ID)
	// подтверждаем отправку сообщения
	n.safeAck(false, delivery)

}
