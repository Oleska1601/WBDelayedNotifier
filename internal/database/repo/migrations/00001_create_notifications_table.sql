-- +goose Up
CREATE TYPE status AS ENUM
(
    'scheduled',    --создано, время еще не наступило
    'sent',         --отправлено
    'cancelled',    --отменено до отправки
    'failed'       --не удалось отправить после всех попыток
);

CREATE TYPE channel AS ENUM
(
    'telegram',    
    'email'         
);

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    channel channel NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    message VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,                   --во сколько было создано
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,                 --на сколько была запланирована отправка
    sent_at TIMESTAMP WITH TIME ZONE,                               --во сколько реально отправилось (учитывая задержки и возможные retry)
    status status NOT NULL 
);


-- +goose Down
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS channel;
DROP TYPE IF EXISTS status;