package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Oleska1601/WBDelayedNotifier/internal/errs"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func (s *Service) GetNotificationStatus(ctx context.Context, id int) (models.Status, error) {
	idStr := strconv.Itoa(id)
	cachedStatus, err := s.cache.GetValue(ctx, idStr)
	if err != nil && !errors.Is(err, redis.NoMatches) {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "service GetNotificationStatus s.cache.GetValue").
			Int("id", id).
			Msg("failed to get cache value")
	}

	if cachedStatus != "" {
		currentStatus := models.Status(cachedStatus)
		if models.IsValidStatus(currentStatus) {
			zlog.Logger.Debug().
				Str("path", "service GetNotificationStatus").
				Int("id", id).
				Msg("get cache value")
			return currentStatus, nil
		}
	}

	status, err := s.repo.GetNotificationStatus(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get notification status: %w", err)
	}

	if err := s.cache.SetValue(ctx, idStr, string(status)); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "service GetNotificationStatus s.cache.SetValue").
			Int("id", id).
			Msg("failed to set cache value")
	}

	return status, nil

}

func (s *Service) CreateNotification(ctx context.Context, notification *models.Notification) (int, error) {
	id, err := s.repo.CreateNotification(ctx, notification)
	if err != nil {
		return 0, fmt.Errorf("create notification: %w", err)
	}

	notification.ID = id
	if err := s.publisher.PublishNotification(ctx, notification); err != nil {
		updateErr := s.updateNotificationStatus(ctx, id, models.StatusFailed)
		if updateErr != nil {
			return 0, fmt.Errorf("publish notification: %w update notification status: %w", err, updateErr)
		}

		return 0, fmt.Errorf("publish notification: %w", err)
	}

	idStr := strconv.Itoa(id)
	if err := s.cache.SetValue(ctx, idStr, string(models.StatusScheduled)); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "service CreateNotification s.cache.SetValue").
			Int("id", id).
			Msg("failed to set cache value")
	}

	return id, nil
}

func (s *Service) DeleteNotification(ctx context.Context, id int) error {
	getStatus, err := s.GetNotificationStatus(ctx, id)
	if err != nil {
		return err
	}

	// processed/cancelled/failed удалению не подлежит
	if getStatus != models.StatusScheduled {
		return errs.NewConflictError(fmt.Sprintf("%s notification cannot be processed", getStatus))
	}

	return s.updateNotificationStatus(ctx, id, models.StatusCancelled)
}

func (s *Service) updateNotificationStatus(ctx context.Context, id int, status models.Status) error {
	err := s.repo.UpdateNotificationStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("update notification status: %w", err)
	}

	idStr := strconv.Itoa(id)
	if err := s.cache.SetValue(ctx, idStr, string(status)); err != nil {
		zlog.Logger.Warn().
			Err(err).
			Str("path", "service updateNotificationStatus s.cache.SetValue").
			Int("id", id).
			Msg("failed to set cache value")
	}

	return nil
}
