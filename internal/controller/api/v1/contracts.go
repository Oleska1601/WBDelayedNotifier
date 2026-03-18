package v1

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type ServiceI interface {
	GetNotificationStatus(context.Context, int) (models.Status, error)
	CreateNotification(context.Context, *models.Notification) (int, error)
	DeleteNotification(context.Context, int) error
}
