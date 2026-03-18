-- +goose Up
CREATE TYPE status AS ENUM
(
    'scheduled',    --создано, время еще не наступило
    'sent',         --отправлено
    'cancelled',    --отменено до отправки
    'failed'        --не удалось отправить после всех попыток
);

CREATE TYPE channel AS ENUM
(
    'telegram',    
    'email'         
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    channel channel NOT NULL,
    recipient VARCHAR(255) NOT NULL,                                -- tgID или email
    message VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),     --во сколько было создано
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,                 --на сколько была запланирована отправка
    sent_at TIMESTAMP WITH TIME ZONE,                               --во сколько реально отправилось (учитывая задержки и возможные retry)
    status status NOT NULL DEFAULT 'scheduled'
);

CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_sent_at_notification()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'sent' AND OLD.status != 'sent' THEN
        NEW.sent_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_sent_notification
  BEFORE UPDATE ON notifications
  FOR EACH ROW
  EXECUTE FUNCTION update_sent_at_notification();

-- +goose Down
DROP TRIGGER IF EXISTS trg_sent_notification ON notifications;
DROP FUNCTION IF EXISTS update_sent_at_notification;

DROP INDEX IF EXISTS idx_notifications_status;

DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS channel;
DROP TYPE IF EXISTS status;