package repo

import (
	"context"
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

const (
	getNotificationStatus = `SELECT status FROM notifications WHERE id = $1`
	insertNotification    = `
		INSERT INTO notifications (
			channel, recipient, message, created_at, scheduled_at, status) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`
	updateNotificationStatus = `UPDATE notifications SET status = $1 WHERE id = $2`
)

func (r *PgRepo) GetNotificationStatus(ctx context.Context, notificationID int64) (models.Status, error) {
	res := r.db.Master.QueryRowContext(ctx, getNotificationStatus, notificationID)
	var status models.Status
	if err := res.Scan(&status); err != nil {
		return "", fmt.Errorf("fmt.Scan: %w", err)
	}
	return status, nil
}

func (r *PgRepo) UpdateReadyNotifications(ctx context.Context) ([]models.Notification, error) {
	res, err := r.db.Master.QueryContext(ctx, getNotificationStatus)
	if err != nil {
		return nil, fmt.Errorf("r.db.Master.QueryContext: %w", err)
	}
	notifications := []models.Notification{}
	var notification models.Notification
	for res.Next() {
		if err := res.Scan(&notification); err != nil {
			return nil, fmt.Errorf("fmt.Scan: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *PgRepo) CreateNotification(ctx context.Context, notification models.Notification) (int64, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}
	defer r.rollbackTransaction(tx)

	res, err := tx.ExecContext(ctx, insertNotification,
		notification.Channel,
		notification.Recipient,
		notification.Message,
		notification.CreatedAt,
		notification.ScheduledAt,
		notification.Status,
	)
	if err != nil {
		return 0, fmt.Errorf("tx.ExecContext: %w", err)
	}
	notificationID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("res.LastInsertId: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx.Commit: %w", err)
	}

	return notificationID, nil
}

func (r *PgRepo) UpdateNotificationStatus(ctx context.Context, notificationID int64, status models.Status) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}
	defer r.rollbackTransaction(tx)

	_, err = tx.ExecContext(ctx, updateNotificationStatus, notificationID, status)
	if err != nil {
		return fmt.Errorf("tx.ExecContext: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil

}
