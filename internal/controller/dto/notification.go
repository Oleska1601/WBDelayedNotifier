package dto

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

type GetNotificationResponse struct {
	NotificationStatus models.Status `json:"notification_status"`
}

type CreateNotificationRequest struct {
	Channel     models.Channel `json:"channel" binding:"required,oneof=email telegram" example:"telegram"`
	Recipient   string         `json:"recipient" binding:"required" example:"796744183"` // адрес тг или email
	Message     string         `json:"message" validate:"required" example:"test_message"`
	ScheduledAt time.Time      `json:"scheduled_at" validate:"required" example:"2025-11-01T10:40:00Z"`
}

type CreateNotificationResponse struct {
	NotificationID int64 `json:"notification_id"`
}

// validation
func (r *CreateNotificationRequest) Validate(createdAt time.Time) error {
	switch r.Channel {
	case models.ChannelEmail:
		if !isValidEmail(r.Recipient) {
			return errors.New("invalid type of email")
		}
	case models.ChannelTelegram:
		if !isValidTelegram(r.Recipient) {
			return errors.New("invalid type of telegram id")
		}
	default:
		return errors.New("type of channel is not supported")
	}

	if strings.TrimSpace(r.Recipient) == "" {
		return errors.New("recipient cannot be empty")
	}

	if r.ScheduledAt.Before(createdAt) {
		return errors.New("scheduled_at cannot be in the past")
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func isValidTelegram(tg string) bool {
	// проверяем числовой chat_id
	chatID, err := strconv.ParseInt(tg, 10, 64)
	if err != nil {
		return false
	}

	return chatID > 0
}

func (r *CreateNotificationRequest) ToModel() (models.Notification, error) {
	createdAt := time.Now()

	if err := r.Validate(createdAt); err != nil {
		return models.Notification{}, err
	}

	status := models.StatusScheduled

	notification := models.Notification{
		Channel:     r.Channel,
		Recipient:   r.Recipient,
		Message:     r.Message,
		CreatedAt:   createdAt,
		ScheduledAt: r.ScheduledAt,
		Status:      status,
	}

	return notification, nil
}
