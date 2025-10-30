package notifier

import (
	"context"
	"encoding/json"
	"sync"

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
		return
	}

	var notification models.NotificationMessage
	// 1
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "process json.Unmarshal").
			Msg("failed to unmarshal json")

		delivery.Nack(false, false)
		return
	}

	// 2
	currentStatus, err := n.usecase.GetNotificationStatus(ctx, notification.ID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "processConsume n.usecase.GetNotificationStatus").
			Msg("failed to get notification status")
		delivery.Nack(false, false)
		return
	}

	// 3
	if currentStatus == models.StatusCancelled {
		zlog.Logger.Info().
			Str("path", "processConsume").
			Int64("notification_id", notification.ID).
			Msg("send notification is cancelled")
		delivery.Ack(false) // сообщение отменено - доставлять никуда не надо -> его обработка завершена
		return
	}

	// 4
	// пробуем отправить по указанному источнику (тг/почта)
	// если не удалось - отправляем на обработку еще раз с заданной экспоненциальной задержкой
	// такой цикл может макс повторятся 5 раз (см RabbitMQConfig)
	if err := sender.Send(notification.Recipient, notification.Message); err != nil {
		n.processRetry(ctx, delivery, notification.ID, channelType)
		return
	}

	err = n.updateDB(ctx, notification.ID, models.StatusSent)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "processConsume n.updateDB").
			Int64("notification_id", notification.ID)
		// мб кинуть алерт, тк состояние бд неконсистентно

	}
	// подтверждаем отправку сообщения
	delivery.Ack(false)

}
