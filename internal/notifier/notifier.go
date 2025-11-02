package notifier

import (
	"context"
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type Notifier struct {
	channel *amqp091.Channel
	cfg     *config.RabbitMQConfig
	usecase NotifyUsecaseInterface
	senders map[models.Channel]SenderInterface
}

func New(cfg *config.RabbitMQConfig, usecase NotifyUsecaseInterface) (*Notifier, error) {
	conn, err := rabbitmq.Connect(cfg.Conn.URL, cfg.Conn.Retries, cfg.Conn.Pause)
	if err != nil {
		return nil, fmt.Errorf("establish a connection to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	props := make(map[string]interface{})
	props["x-delayed-type"] = cfg.RetryExchange.Type
	// Declare the exchange
	err = channel.ExchangeDeclare(
		cfg.RetryExchange.Name,
		cfg.RetryExchange.Type,
		cfg.RetryExchange.Durable,
		cfg.RetryExchange.AutoDelete,
		cfg.RetryExchange.Internal,
		cfg.RetryExchange.NoWait,
		props,
	)
	if err != nil {
		return nil, fmt.Errorf("exchange declare: %w", err)
	}
	_, err = channel.QueueDeclare(
		cfg.QueueEmail.Name,
		cfg.QueueEmail.Durable,
		cfg.QueueEmail.AutoDelete,
		cfg.QueueEmail.Exclusive,
		cfg.QueueEmail.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("email queue declare: %w", err)
	}

	err = channel.QueueBind(
		cfg.QueueEmail.Name,
		cfg.QueueEmail.RoutingKey,
		cfg.Exchange.Name,
		cfg.QueueEmail.NoWait,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("email queue bind: %w", err)
	}
	_, err = channel.QueueDeclare(
		cfg.QueueTg.Name,
		cfg.QueueTg.Durable,
		cfg.QueueTg.AutoDelete,
		cfg.QueueTg.Exclusive,
		cfg.QueueTg.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("tg queue declare: %w", err)
	}

	err = channel.QueueBind(
		cfg.QueueTg.Name,
		cfg.QueueTg.RoutingKey,
		cfg.Exchange.Name,
		cfg.QueueTg.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("tg queue bind: %w", err)
	}

	// 3. Retry очереди
	_, err = channel.QueueDeclare(
		cfg.RetryQueueEmail.Name,
		cfg.RetryQueueEmail.Durable,    // durable
		cfg.RetryQueueEmail.AutoDelete, // auto-delete
		cfg.RetryQueueEmail.Exclusive,  // exclusive
		cfg.RetryQueueEmail.NoWait,     // no-wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("email retry queue declare: %w", err)
	}

	_, err = channel.QueueDeclare(
		cfg.RetryQueueTg.Name,
		cfg.RetryQueueTg.Durable,    // durable
		cfg.RetryQueueTg.AutoDelete, // auto-delete
		cfg.RetryQueueTg.Exclusive,  // exclusive
		cfg.RetryQueueTg.NoWait,     // no-wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("tg retry queue declare: %w", err)
	}

	// 5. Биндинги для retry очередей
	err = channel.QueueBind(
		cfg.RetryQueueEmail.Name,
		cfg.RetryQueueEmail.RoutingKey,
		cfg.RetryExchange.Name,
		false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("retry email queue bind: %w", err)
	}

	err = channel.QueueBind(
		cfg.RetryQueueTg.Name,
		cfg.RetryQueueTg.RoutingKey,
		cfg.RetryExchange.Name,
		false, nil,
	)

	if err != nil {
		return nil, fmt.Errorf("retry tg queue bind: %w", err)
	}

	return &Notifier{
		channel: channel,
		cfg:     cfg,
		usecase: usecase,
		senders: make(map[models.Channel]SenderInterface),
	}, nil

}

func (n *Notifier) RegisterSender(channel models.Channel, sender SenderInterface) {
	n.senders[channel] = sender
}

func (n *Notifier) StartWorkers(ctx context.Context) {
	// Основные воркеры - высокий приоритет
	go n.consume(ctx, workers, n.cfg.QueueEmail.Name, models.ChannelEmail)
	go n.consume(ctx, workers, n.cfg.QueueTg.Name, models.ChannelTelegram)

	// Воркеры для ретраев - низкий приоритет, можно меньше ресурсов
	go n.consume(ctx, retryWorkers, n.cfg.RetryQueueEmail.Name, models.ChannelEmail)
	go n.consume(ctx, retryWorkers, n.cfg.RetryQueueTg.Name, models.ChannelTelegram)
}

func (n *Notifier) safeAck(multiple bool, delivery amqp091.Delivery) {
	if err := delivery.Ack(multiple); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message_id", delivery.MessageId).
			Msg("failed to ack message")
	}
}

func (n *Notifier) safeNack(multiple bool, requeue bool, delivery amqp091.Delivery) {
	if err := delivery.Nack(multiple, requeue); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message_id", delivery.MessageId).
			Msg("failed to nack message")
	}
}
