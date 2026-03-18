package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/internal/errs"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
)

const (
	selectNotificationStatus = `SELECT status FROM notifications WHERE id = $1`
	insertNotification       = `INSERT INTO notifications (channel, recipient, message, scheduled_at) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	updateNotificationStatus = `UPDATE notifications SET status = $1 WHERE id = $2`
)

func (r *PgRepo) GetNotificationStatus(ctx context.Context, id int) (models.Status, error) {
	row := r.db.Master.QueryRowContext(ctx, selectNotificationStatus, id)
	var status models.Status
	if err := row.Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NewNotFoundError(fmt.Sprintf("get notification status with id %d", id))
		}

		return "", fmt.Errorf("row scan: %w", err)
	}

	return status, nil
}

func (r *PgRepo) CreateNotification(ctx context.Context, notification *models.Notification) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, insertNotification,
		notification.Channel,
		notification.Recipient,
		notification.Message,
		notification.ScheduledAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("insert notification: %w", err)
	}

	return id, nil
}

func (r *PgRepo) UpdateNotificationStatus(ctx context.Context, id int, status models.Status) error {
	res, err := r.db.ExecContext(ctx, updateNotificationStatus, status, id)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rows == 0 {
		return errs.NewNotFoundError(fmt.Sprintf("update notification status with id %d", id))
	}

	return nil
}
