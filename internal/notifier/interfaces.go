package notifier

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type SenderInterface interface {
	Send(recipient, message string) error // parsing
}

type NotifyUsecaseInterface interface {
	GetNotificationStatus(context.Context, int64) (models.Status, error)
	UpdateNotification(context.Context, models.UpdateNotification) error
}
