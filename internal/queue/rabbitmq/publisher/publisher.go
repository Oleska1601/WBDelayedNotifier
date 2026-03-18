package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	rmq "github.com/Oleska1601/WBDelayedNotifier/internal/queue/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

const (
	attempts = 3
	delay    = time.Millisecond * 100
	backoff  = 2
)

type Publisher struct {
	publisher *rabbitmq.Publisher
}

func New(cfg *config.RabbitMQConfig) (*Publisher, error) {
	clientCfg := rabbitmq.ClientConfig{
		URL:            cfg.ClientConn.URL,
		ProducingStrat: retry.Strategy{Attempts: attempts, Delay: delay, Backoff: backoff},
	}
	client, err := rabbitmq.NewClient(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("create new client: %w", err)
	}

	args := make(map[string]interface{})
	args["x-delayed-type"] = rmq.ExchangeType
	err = client.DeclareExchange(
		rmq.ExchangeName,
		rmq.ExchangeKind,
		rmq.ExchangeDurable,
		rmq.ExchangeAutoDelete,
		rmq.ExchangeInternal,
		args,
	)
	if err != nil {
		return nil, fmt.Errorf("declare declare: %w", err)
	}

	err = client.DeclareExchange(
		rmq.DeadLetterExchangeName,
		rmq.DeadLetterExchangeKind,
		rmq.DeadLetterExchangeDurable,
		rmq.DeadLetterExchangeAutoDelete,
		rmq.DeadLetterExchangeInternal,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("declare dead letter declare: %w", err)
	}

	dlxArgs := amqp091.Table{
		"x-dead-letter-exchange": rmq.DeadLetterExchangeName,
	}
	err = client.DeclareQueue(
		rmq.EmailQueueName,
		rmq.ExchangeName,
		rmq.EmailRoutingKey,
		rmq.QueueDurable,
		rmq.QueueAutoDelete,
		rmq.ExchangeDurable,
		dlxArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("email queue declare: %w", err)
	}

	err = client.DeclareQueue(
		rmq.TgQueueName,
		rmq.ExchangeName,
		rmq.TgRoutingKey,
		rmq.QueueDurable,
		rmq.QueueAutoDelete,
		rmq.ExchangeDurable,
		dlxArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("tg queue declare: %w", err)
	}

	err = client.DeclareQueue(
		rmq.DeadLetterQueueName,
		rmq.DeadLetterExchangeName,
		rmq.DeadLetterRoutingKey,
		rmq.DeadLetterQueueDurable,
		rmq.DeadLetterQueueAutoDelete,
		rmq.DeadLetterExchangeDurable,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("dead letter queue declare: %w", err)
	}

	publisher := rabbitmq.NewPublisher(client, rmq.ExchangeName, rmq.ContentType)
	return &Publisher{publisher: publisher}, nil
}

func (p *Publisher) PublishNotification(ctx context.Context, notification *models.Notification) error {
	body, err := json.Marshal(*notification)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	var routingKey string
	if notification.Channel == models.ChannelEmail {
		routingKey = rmq.EmailRoutingKey
	} else {
		routingKey = rmq.TgRoutingKey
	}

	delay := notification.ScheduledAt.Sub(time.Now())
	if delay < 0 {
		delay = 0
	}

	err = p.publisher.Publish(ctx, body, routingKey, rabbitmq.WithHeaders(
		amqp091.Table{"x-delay": delay.Milliseconds()},
	))
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
