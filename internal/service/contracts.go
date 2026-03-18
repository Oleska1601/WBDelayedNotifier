package service

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type CacheI interface {
	GetValue(context.Context, string) (string, error)
	SetValue(context.Context, string, interface{}) error
}

type RepoI interface {
	GetNotificationStatus(context.Context, int) (models.Status, error)
	CreateNotification(context.Context, *models.Notification) (int, error)
	UpdateNotificationStatus(context.Context, int, models.Status) error
}

type PublisherI interface {
	PublishNotification(context.Context, *models.Notification) error
}
