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
	updateNotification = `UPDATE notifications SET status = $1, sent_at = $2 WHERE id = $3`
)

func (r *PgRepo) GetNotificationStatus(ctx context.Context, notificationID int64) (models.Status, error) {
	res := r.db.Master.QueryRowContext(ctx, getNotificationStatus, notificationID)
	var status models.Status
	if err := res.Scan(&status); err != nil {
		return "", fmt.Errorf("fmt.Scan: %w", err)
	}
	return status, nil
}

func (r *PgRepo) CreateNotification(ctx context.Context, notification models.Notification) (int64, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}
	defer r.rollbackTransaction(tx)

	var notificationID int64

	err = tx.QueryRowContext(ctx, insertNotification,
		notification.Channel,
		notification.Recipient,
		notification.Message,
		notification.CreatedAt,
		notification.ScheduledAt,
		notification.Status,
	).Scan(&notificationID)

	if err != nil {
		return 0, fmt.Errorf("tx.QueryRowContext: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("tx.Commit: %w", err)
	}

	return notificationID, nil
}

func (r *PgRepo) UpdateNotification(ctx context.Context, notification models.UpdateNotification) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("r.db.Master.BeginTx: %w", err)
	}
	defer r.rollbackTransaction(tx)

	_, err = tx.ExecContext(ctx, updateNotification, notification.Status, notification.SentAt, notification.ID)
	if err != nil {
		return fmt.Errorf("tx.ExecContext: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil

}
