-- +goose Up
BEGIN;

CREATE INDEX IF NOT EXISTS idx_notifications_send_at ON notifications(sent_at);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);

COMMIT;

-- +goose Down
BEGIN;

DROP INDEX IF EXISTS idx_notifications_send_at;
DROP INDEX IF EXISTS idx_notifications_status;

COMMIT;