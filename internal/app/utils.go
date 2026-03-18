package app

import (
	"fmt"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/wb-go/wbf/dbpg"
)

func initDB(cfg *config.PostgresConfig) (*dbpg.DB, error) {
	masterDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	slavesDSN := []string{}
	options := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}
	db, err := dbpg.New(masterDSN, slavesDSN, options)
	if err != nil {
		return nil, fmt.Errorf("create new DB instance: %w", err)
	}

	return db, nil
}

func initTgAPI(cfg *config.TgBotConfig) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("new bot api: %w", err)
	}

	return bot, nil
}
