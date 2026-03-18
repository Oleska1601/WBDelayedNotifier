package models

import "time"

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusSent      Status = "sent"
	StatusCancelled Status = "cancelled"
	StatusFailed    Status = "failed"
)

func IsValidStatus(currentStatus Status) bool {
	return currentStatus == StatusScheduled ||
		currentStatus == StatusSent ||
		currentStatus == StatusCancelled ||
		currentStatus == StatusFailed
}

type Channel string

const (
	ChannelTelegram Channel = "telegram"
	ChannelEmail    Channel = "email"
)

type Notification struct {
	ID          int
	Channel     Channel
	Recipient   string
	Message     string
	CreatedAt   time.Time
	ScheduledAt time.Time
	SentAt      time.Time
	Status      Status
}
