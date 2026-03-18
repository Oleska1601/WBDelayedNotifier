package tgbot

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/wb-go/wbf/zlog"
)

const (
	timeout = 60
)

type TGBotSender struct {
	botAPI *tgbotapi.BotAPI
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func New(botAPI *tgbotapi.BotAPI) *TGBotSender {
	return &TGBotSender{
		botAPI: botAPI,
	}
}

func (tg *TGBotSender) Start(ctx context.Context) {
	tgCtx, cancel := context.WithCancel(ctx)
	tg.cancel = cancel

	tg.wg.Add(1)
	go func() {
		defer tg.wg.Done()
		tg.start(tgCtx)
	}()
}

func (tg *TGBotSender) start(ctx context.Context) {
	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = timeout
	updates, err := tg.botAPI.GetUpdatesChan(updateCfg)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "tgbot Start tg.botAPI.GetUpdatesChan").
			Msg("failed to get channel of updates")
		return
	}

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
				msg := fmt.Sprintf("Hello! Your chatID is %d, please use it for testing getting notifications", chatID)
				if err := tg.SendMessage(chatID, msg); err != nil {
					zlog.Logger.Error().
						Err(err).
						Str("path", "tgbot Start tg.SendMessage").
						Int64("chat_id", chatID).
						Msg("failed to send message")
				}
			}
		}
	}
}

func (tg *TGBotSender) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := tg.botAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (tg *TGBotSender) Stop() {
	tg.cancel()
	tg.wg.Wait()
}
