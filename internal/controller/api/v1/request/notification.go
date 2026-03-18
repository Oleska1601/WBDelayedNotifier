package request

import (
	"errors"
	"net/mail"
	"strconv"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type CreateNotificationRequest struct {
	Channel     models.Channel `json:"channel" binding:"required"`
	Recipient   string         `json:"recipient" binding:"required"` // тг или email
	Message     string         `json:"message" binding:"required"`
	ScheduledAt time.Time      `json:"scheduled_at" binding:"required"`
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidTelegram(chatIDStr string) bool {
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return false
	}

	return chatID > 0
}

func (r *CreateNotificationRequest) ToModel() (*models.Notification, error) {
	switch r.Channel {
	case models.ChannelEmail:
		if !isValidEmail(r.Recipient) {
			return nil, errors.New("invalid type of email")
		}
	case models.ChannelTelegram:
		if !isValidTelegram(r.Recipient) {
			return nil, errors.New("invalid type of telegram id")
		}
	default:
		return nil, errors.New("type of channel is not supported")
	}

	if r.ScheduledAt.Before(time.Now()) {
		return nil, errors.New("scheduled_at cannot be in the past")
	}

	return &models.Notification{
		Channel:     r.Channel,
		Recipient:   r.Recipient,
		Message:     r.Message,
		ScheduledAt: r.ScheduledAt,
	}, nil
}
