package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/zlog"
)

func (u *Usecase) GetNotificationStatus(ctx context.Context, notificationID int64) (models.Status, error) {
	strNotificationID := strconv.FormatInt(notificationID, 10)
	statusStr, err := u.cache.GetValue(ctx, strNotificationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		zlog.Logger.Error().
			Err(err).
			Str("message", "GetNotificationStatus u.cache.GetValue")
		// даже если могла возникнуть ошибка с получением значения из кеша, все равно смотрим в ьд
	}

	if statusStr != "" {
		return models.Status(statusStr), nil
	}

	status, err := u.repo.GetNotificationStatus(ctx, notificationID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "GetNotificationStatus u.repo.GetNotificationStatus")
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("notification with provided notification_id does not exist")
		}
		return "", errors.New("failed to get notification status")
	}
	if err := u.cache.SetValue(ctx, strNotificationID, status); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "GetNotificationStatus u.cache.SetValue")
	}

	return status, nil

}

func (u *Usecase) CreateNotification(ctx context.Context, notification models.Notification) (int64, error) {
	notificationID, err := u.repo.CreateNotification(ctx, notification)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "CreateNotification u.CreateNotification")
		return 0, errors.New("failed to create notification")
	}
	if err := u.publisher.PublishNotification(notification); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "u.publisher.PublishNotification")
		return 0, errors.New("failed to publish notification")
	}
	return notificationID, nil
}

func (u *Usecase) UpdateNotificationStatus(ctx context.Context, notificationID int64, status models.Status) error {
	err := u.repo.UpdateNotificationStatus(ctx, notificationID, status)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "UpdateNotificationStatus u.repo.UpdateNotificationStatus")
		return errors.New("failed to update notification status")
	}

	strNotificationID := strconv.FormatInt(notificationID, 10)
	if err := u.cache.SetValue(ctx, strNotificationID, status); err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("message", "UpdateNotificationStatus u.cache.SetValue")
	}
	// даже если ошибка в кешэ - все равно возращаем nil, тк значение получить удалось
	return nil
}
