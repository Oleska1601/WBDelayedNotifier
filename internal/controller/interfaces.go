package controller

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type UsecaseInterface interface {
	GetNotificationStatus(context.Context, int64) (models.Status, error)
	CreateNotification(context.Context, models.Notification) (int64, error)
	UpdateNotificationStatus(context.Context, int64, models.Status) error
}
