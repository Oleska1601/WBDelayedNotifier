package consumer

import (
	"context"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type RepoI interface {
	GetNotificationStatus(context.Context, int) (models.Status, error)
	UpdateNotificationStatus(context.Context, int, models.Status) error
}

type EmailI interface {
	SendMessage(string, string) error
}

type TgBotI interface {
	SendMessage(int64, string) error
}
