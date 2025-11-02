package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/controller"
	"github.com/Oleska1601/WBDelayedNotifier/internal/database/repo"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/Oleska1601/WBDelayedNotifier/internal/notifier"
	"github.com/Oleska1601/WBDelayedNotifier/internal/publisher"
	"github.com/Oleska1601/WBDelayedNotifier/internal/redis"
	"github.com/Oleska1601/WBDelayedNotifier/internal/sender/email"
	"github.com/Oleska1601/WBDelayedNotifier/internal/sender/tgbot"
	"github.com/Oleska1601/WBDelayedNotifier/internal/usecase"
	"github.com/wb-go/wbf/zlog"
)

// @title Delayed Notifier
// @version 1.0
// @description API for Delayed Notifier
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /
func Run(cfg *config.Config) {

	// logger
	zlog.Init()
	if err := zlog.SetLevel(cfg.Logger.Level); err != nil {
		log.Fatalln("set zlog level error: %w", err)
	}

	// postgres
	db, err := initDB(&cfg.DB)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run initDB").
			Msg("init database")
	}

	// postgres repo
	pgRepo := repo.New(db)
	if err := pgRepo.ApplyMigrations(); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run pgRepo.ApplyMigrations").
			Msg("apply migrations to database")
	}
	redis, err := redis.New(&cfg.Redis)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run redis.New").
			Msg("init redis")
	}

	// publisher
	publisher, err := publisher.New(&cfg.RabbitMQ)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run publisher.New").
			Msg("init publisher")
	}

	usecase := usecase.New(redis, pgRepo, publisher)
	server := controller.New(&cfg.Server, usecase)

	// tgbot sender
	tgbotSender, err := tgbot.New(&cfg.Telegram)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run tgbot.New").
			Msg("init tgbot")
	}

	// email sender
	emailSender := email.New(&cfg.Email)

	notifier, err := notifier.New(&cfg.RabbitMQ, usecase)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run notifier.New").
			Msg("init notifier")
	}
	notifier.RegisterSender(models.ChannelEmail, emailSender)
	notifier.RegisterSender(models.ChannelTelegram, tgbotSender)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go tgbotSender.HandleUpdates(ctx)

	notifier.StartWorkers(ctx)

	go func() {
		if err := server.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.
				Fatal().
				Err(err).
				Str("path", "Run server.Srv.ListenAndServe").
				Msg("cannot start server")
		}
		zlog.Logger.Info().Msgf("server is started http://%s:%d/", cfg.Server.Host, cfg.Server.Port)
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Logger.Info().Msg("shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Srv.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("server.Srv.Shutdown")
		return
	}

	zlog.Logger.Info().Msg("server exited properly")

}
