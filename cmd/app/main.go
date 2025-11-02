package main

import (
	"log/slog"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("main config.New", slog.Any("error", err))
		return
	}
	app.Run(cfg)

}
