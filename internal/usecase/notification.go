package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/zlog"
)

func (u *Usecase) GetNotificationStatus(ctx context.Context, notificationID int64) (models.Status, error) {
	notificationIDStr := strconv.FormatInt(notificationID, 10)
	statusStr, err := u.cache.GetValue(ctx, notificationIDStr)
	if err != nil && !errors.Is(err, redis.Nil) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetNotificationStatus u.cache.GetValue").
			Int64("notification_id", notificationID).
			Msg("failed to get cache value")
	}

	if statusStr != "" {
		currentStatus := models.Status(statusStr)
		if models.IsValidStatus(currentStatus) {
			return currentStatus, nil
		}
	}

	status, err := u.repo.GetNotificationStatus(ctx, notificationID)
	if err != nil {
		return "", fmt.Errorf("GetNotificationStatus u.repo.GetNotificationStatus: %w", err)
	}

	if err := u.cache.SetValue(ctx, notificationIDStr, string(status)); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "GetNotificationStatus u.cache.SetValue").
			Int64("notification_id", notificationID).
			Msg("failed to set cache value")
	}

	return status, nil

}

func (u *Usecase) CreateNotification(ctx context.Context, notification models.Notification) (int64, error) {
	notificationID, err := u.repo.CreateNotification(ctx, notification)
	if err != nil {
		return 0, fmt.Errorf("CreateNotification u.repo.CreateNotification: %w", err)
	}

	notification.ID = notificationID
	if err := u.publisher.PublishNotification(notification); err != nil {
		updateNotification := models.UpdateNotification{
			ID:     notificationID,
			Status: models.StatusFailed,
		}
		updateErr := u.UpdateNotification(ctx, updateNotification)
		if updateErr != nil {
			return 0, fmt.Errorf("u.publisher.PublishNotification: %w u.UpdateNotification: %w", err, updateErr)
		}

		return 0, fmt.Errorf("u.publisher.PublishNotification: %w", err)
	}

	notificationIDStr := strconv.FormatInt(notificationID, 10)

	if err := u.cache.SetValue(ctx, notificationIDStr, string(notification.Status)); err != nil { // notification.Status: publsihed
		zlog.Logger.Warn().
			Err(err).
			Str("path", "CreateNotification u.cache.SetValue").
			Int64("notification_id", notificationID).
			Msg("failed to set cache value")
	}

	return notificationID, nil
}

func (u *Usecase) UpdateNotification(ctx context.Context, updateNotification models.UpdateNotification) error {
	// проверяем сначала текущ статус и вообще есть ли такой id
	getStatus, err := u.GetNotificationStatus(ctx, updateNotification.ID)
	if err != nil {
		return fmt.Errorf("u.GetNotificationStatus: %w", err)
	}

	if getStatus == models.StatusCancelled {
		return errors.New("notification is already cancelled")
	}

	err = u.repo.UpdateNotification(ctx, updateNotification)
	if err != nil {
		return fmt.Errorf("u.repo.UpdateNotification: %w", err)
	}

	notificationIDStr := strconv.FormatInt(updateNotification.ID, 10)
	if err := u.cache.SetValue(ctx, notificationIDStr, string(updateNotification.Status)); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "UpdateNotification u.cache.SetValue").
			Int64("notification_id", updateNotification.ID).
			Msg("failed to set cache value")
	}

	return nil
}
