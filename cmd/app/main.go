package main

import (
	"log/slog"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/app"
)

// @title Delayed Notifier
// @version 1.0
// @description API for Order Service
// @termsOfService http://swagger.io/terms/

// @host :8081
// @BasePath /
func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("main config.New", slog.Any("error", err))
		return
	}
	app.Run(cfg)

}
