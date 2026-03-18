package email

import (
	"fmt"
	"net/smtp"

	"github.com/Oleska1601/WBDelayedNotifier/config"
)

type EmailSender struct {
	cfg *config.EmailConfig
}

func New(cfg *config.EmailConfig) *EmailSender {
	return &EmailSender{
		cfg: cfg,
	}
}

func (es *EmailSender) SendMessage(recipient, message string) error {
	auth := smtp.PlainAuth("", es.cfg.Username, es.cfg.Password, es.cfg.Host)
	addr := fmt.Sprintf("%s:%d", es.cfg.Host, es.cfg.Port)
	err := smtp.SendMail(addr, auth, es.cfg.From, []string{recipient}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
