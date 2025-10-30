package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/wb-go/wbf/rabbitmq"
)

type Publisher struct {
	publisher *rabbitmq.Publisher
	//channel   *ampq091.Channel
	cfg config.RabbitMQConfig
}

func New(cfg config.RabbitMQConfig) (*Publisher, error) {
	conn, err := rabbitmq.Connect(cfg.Conn.URL, cfg.Conn.Retries, cfg.Conn.Pause)
	if err != nil {
		return nil, fmt.Errorf("establish a connection to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	props := make(map[string]interface{})
	props["x-delayed-type"] = cfg.Exchange.Type
	// Declare the exchange
	err = channel.ExchangeDeclare(
		cfg.Exchange.Name,
		cfg.Exchange.Type,
		cfg.Exchange.Durable,
		cfg.Exchange.AutoDelete,
		cfg.Exchange.Internal,
		cfg.Exchange.NoWait,
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

	publisher := rabbitmq.NewPublisher(channel, cfg.Exchange.Name)
	return &Publisher{
		publisher: publisher,
		cfg:       cfg,
	}, nil

}

func (p *Publisher) PublishNotification(notification models.Notification) error {

	// Set the message delay time using the message header, in milliseconds
	// Вычисляем задержку в миллисекундах
	createdAt := notification.CreatedAt
	scheduledAt := notification.ScheduledAt
	delay := scheduledAt.Sub(createdAt)

	if delay < 0 {
		delay = 0

	}

	var routingKey string
	if notification.Channel == "telegram" {
		routingKey = p.cfg.QueueTg.RoutingKey
	} else {
		routingKey = p.cfg.QueueEmail.RoutingKey
	}

	message := models.NotificationMessage{
		ID:        notification.ID,
		Recipient: notification.Recipient,
		Message:   notification.Message,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	msgHeaders := make(map[string]interface{})
	msgHeaders["x-delay"] = delay
	options := rabbitmq.PublishingOptions{
		Headers: msgHeaders,
	}

	if err := p.publisher.Publish(body, routingKey, "application/json", options); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}
