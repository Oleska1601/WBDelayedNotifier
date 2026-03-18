package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	rmq "github.com/Oleska1601/WBDelayedNotifier/internal/queue/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

const (
	attemts = 5
	delay   = time.Millisecond * 100
	backoff = 2

	workers           = 3
	deadLetterWorkers = 2
)

type Consumer struct {
	emailConsumer      *rabbitmq.Consumer
	tgConsumer         *rabbitmq.Consumer
	deadLetterConsumer *rabbitmq.Consumer
	repo               RepoI
	email              EmailI
	tgBot              TgBotI

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func New(cfg *config.RabbitMQConfig, repo RepoI, email EmailI, tgBot TgBotI) (*Consumer, error) {
	strategy := retry.Strategy{Attempts: attemts, Delay: delay, Backoff: backoff}
	clientCfg := rabbitmq.ClientConfig{URL: cfg.ClientConn.URL, ConnectTimeout: cfg.ClientConn.Timeout, ConsumingStrat: strategy}
	client, err := rabbitmq.NewClient(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("create new client: %w", err)
	}

	emailConsumerCfg := rabbitmq.ConsumerConfig{
		Queue:       rmq.EmailQueueName,
		ConsumerTag: rmq.EmailRoutingKey,
		Workers:     workers,
	}

	tgConsumerCfg := rabbitmq.ConsumerConfig{
		Queue:       rmq.TgQueueName,
		ConsumerTag: rmq.TgRoutingKey,
		Workers:     workers,
	}

	deadLetterConsumerCfg := rabbitmq.ConsumerConfig{
		Queue:       rmq.DeadLetterQueueName,
		ConsumerTag: rmq.DeadLetterRoutingKey,
		Workers:     deadLetterWorkers,
	}

	c := &Consumer{
		repo:  repo,
		email: email,
		tgBot: tgBot,
	}
	emailConsumer := rabbitmq.NewConsumer(client, emailConsumerCfg, c.sendEmailMessage)
	c.emailConsumer = emailConsumer
	tgConsumer := rabbitmq.NewConsumer(client, tgConsumerCfg, c.sendTgMessage)
	c.tgConsumer = tgConsumer
	deadLetterConsumer := rabbitmq.NewConsumer(client, deadLetterConsumerCfg, c.setFailedNotificationStatus)
	c.deadLetterConsumer = deadLetterConsumer
	return c, nil
}

func (c *Consumer) Start(ctx context.Context) {
	consumerCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.wg.Add(3)
	go func() {
		defer c.wg.Done()
		if err := c.emailConsumer.Start(consumerCtx); err != nil {
			zlog.Logger.Error().
				Err(err).
				Str("path", "consumer Start c.emailConsumer.Start").
				Msg("failed to start email consumer")
		}
	}()

	go func() {
		defer c.wg.Done()
		if err := c.tgConsumer.Start(consumerCtx); err != nil {
			zlog.Logger.Error().
				Err(err).
				Str("path", "consumer Start c.tgConsumer.Start").
				Msg("failed to start tg consumer")
		}
	}()

	go func() {
		defer c.wg.Done()
		if err := c.deadLetterConsumer.Start(consumerCtx); err != nil {
			zlog.Logger.Error().
				Err(err).
				Str("path", "consumer Start c.deadLetterConsumer.Start").
				Msg("failed to start dead letter consumer")
		}
	}()
}

func (c *Consumer) sendEmailMessage(ctx context.Context, delivery amqp091.Delivery) error {
	var notification models.Notification
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "sendEmailMessage json.Unmarshal").
			Msg("failed to unmarshal json")
		return nil
	}

	status, err := c.repo.GetNotificationStatus(ctx, notification.ID)
	if err != nil {
		return fmt.Errorf("get notification status: %w", err)
	}

	if status == models.StatusCancelled {
		zlog.Logger.Debug().
			Str("path", "sendEmailMessage").
			Int("id", notification.ID).
			Msg("skipped: status is cancelled")
		return nil
	}

	msg := formNotificationMessage(notification.Message)
	if err := c.email.SendMessage(notification.Recipient, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *Consumer) sendTgMessage(ctx context.Context, delivery amqp091.Delivery) error {
	notification, err := c.processDelivery(ctx, delivery)
	if err != nil {
		return fmt.Errorf("failed to process delivery: %w", err)
	}

	if notification == nil {
		return nil
	}

	chatID, err := strconv.ParseInt(notification.Recipient, 10, 64)
	if err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "sendTgMessage strconv.ParseInt").
			Int("id", notification.ID).
			Msg("failed to parse recipient")
		return nil
	}

	msg := formNotificationMessage(notification.Message)
	if err := c.tgBot.SendMessage(chatID, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *Consumer) setFailedNotificationStatus(ctx context.Context, delivery amqp091.Delivery) error {
	notification, err := c.processDelivery(ctx, delivery)
	if err != nil {
		return fmt.Errorf("failed to process delivery: %w", err)
	}

	if notification == nil {
		return nil
	}

	if err := c.repo.UpdateNotificationStatus(ctx, notification.ID, models.StatusFailed); err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

func (c *Consumer) Stop() {
	c.cancel()
	c.wg.Wait()
}
