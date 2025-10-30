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
	Channel     models.Channel `json:"channel" binding:"required,oneof=email telegram" example:"email"`
	Recipient   string         `json:"recipient" binding:"required" example:"email@gmail.com"` // адрес тг или email
	Message     string         `json:"message" validate:"required" example:"test_message"`
	ScheduledAt time.Time      `json:"scheduled_at" validate:"required" example:"2024-01-15T14:00:00Z"`
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
	// Проверяем либо @username, либо числовой ID
	if strings.HasPrefix(tg, "@") {
		return len(tg) >= 5 && len(tg) <= 32 // @username от 5 до 32 символов
	}
	// Или числовой ID
	_, err := strconv.ParseInt(tg, 10, 64)
	return err == nil
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
