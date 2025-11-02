package usecase

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type CacheInterface interface {
	GetValue(context.Context, string) (string, error)
	SetValue(context.Context, string, any) error
}

type PublisherInterface interface {
	PublishNotification(models.Notification) error
}

type RepoInterface interface {
	GetNotificationStatus(context.Context, int64) (models.Status, error)
	CreateNotification(context.Context, models.Notification) (int64, error)
	UpdateNotification(context.Context, models.UpdateNotification) error
}
