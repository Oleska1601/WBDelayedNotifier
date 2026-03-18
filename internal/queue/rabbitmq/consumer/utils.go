package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/zlog"
)

func formNotificationMessage(message string) string {
	return fmt.Sprintf("Hello from notification service!\n Your notification: %s\n", message)
}

func (c *Consumer) processDelivery(ctx context.Context, delivery amqp091.Delivery) (*models.Notification, error) {
	var notification models.Notification
	if err := json.Unmarshal(delivery.Body, &notification); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "processDelivery json.Unmarshal").
			Msg("failed to unmarshal json")
		return nil, nil
	}

	status, err := c.repo.GetNotificationStatus(ctx, notification.ID)
	if err != nil {
		return nil, fmt.Errorf("get notification status: %w", err)
	}

	if status == models.StatusCancelled {
		zlog.Logger.Debug().
			Str("path", "processDelivery").
			Int("id", notification.ID).
			Msg("skipped: status is cancelled")
		return nil, nil
	}

	return &notification, nil
}
