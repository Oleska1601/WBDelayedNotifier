package tgbot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wb-go/wbf/zlog"
)

const (
	timeout = 60
)

type TGBotSender struct {
	bot *tgbotapi.BotAPI
}

func New(token string) (*TGBotSender, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("creates a new BotAPI instance: %w", err)
	}
	bot.Debug = false

	return &TGBotSender{
		bot: bot,
	}, nil
}

func (tg *TGBotSender) Send(chatIDStr, message string) error {
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("convert chatID to int64: %w", err)
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, err = tg.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

func (ts *TGBotSender) HandleUpdates(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := ts.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}

			if update.Message != nil && update.Message.Text == "/start" {
				chatID := update.Message.Chat.ID
				reply := fmt.Sprintf("Hello! Your chatID is %d, please use it for testing getting notifications", chatID)

				msg := tgbotapi.NewMessage(chatID, reply)
				_, err := ts.bot.Send(msg)
				if err != nil {
					zlog.Logger.Error().
						Err(err).
						Str("path", "HandleUpdates ts.bot.Send").
						Msg("failed to send message")
				}
			}
		}
	}
}
