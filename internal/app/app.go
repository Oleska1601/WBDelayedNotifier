package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oleska1601/WBDelayedNotifier/config"
	"github.com/Oleska1601/WBDelayedNotifier/internal/controller/api"
	v1 "github.com/Oleska1601/WBDelayedNotifier/internal/controller/api/v1"
	"github.com/Oleska1601/WBDelayedNotifier/internal/queue/rabbitmq/consumer"
	"github.com/Oleska1601/WBDelayedNotifier/internal/queue/rabbitmq/publisher"
	"github.com/Oleska1601/WBDelayedNotifier/internal/redis"
	"github.com/Oleska1601/WBDelayedNotifier/internal/repo/postgres"
	"github.com/Oleska1601/WBDelayedNotifier/internal/sender/email"
	"github.com/Oleska1601/WBDelayedNotifier/internal/sender/tgbot"
	"github.com/Oleska1601/WBDelayedNotifier/internal/service"
	"github.com/wb-go/wbf/zlog"
)

// @title Delayed Notifier
// @version 1.0
// @description API for Delayed Notifier
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /
func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zlog.Init()
	if err := zlog.SetLevel(cfg.Logger.Level); err != nil {
		log.Fatalln("set zlog level error: %w", err)
	}

	db, err := initDB(&cfg.Database.Postgres)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run initDB").
			Msg("init database")
	}

	pgRepo := postgres.New(db)
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

	publisher, err := publisher.New(&cfg.RabbitMQ)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run publisher.New").
			Msg("init publisher")
	}
	emailSender := email.New(&cfg.Email)
	tgAPI, err := initTgAPI(&cfg.TgBot)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run initTgAPI").
			Msg("failed to init tg api")
	}

	tgBotSender := tgbot.New(tgAPI)
	consumer, err := consumer.New(&cfg.RabbitMQ, pgRepo, emailSender, tgBotSender)
	service := service.New(redis, pgRepo, publisher)
	apiV1 := v1.New(service)
	router := api.Register(&cfg.Gin, apiV1)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{Addr: addr, Handler: router}

	tgBotSender.Start(ctx)
	consumer.Start(ctx)

	go func() {
		zlog.Logger.Info().Str("path", "Run").Str("addr", addr).Msg("start server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlog.Logger.Fatal().
				Err(err).
				Str("path", "Run server.ListenAndServe").
				Msg("failed to process server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Err(err).Str("path", "App server.Shutdown").
			Msg("failed to shutdown server")
	}

	tgBotSender.Stop()
	consumer.Stop()

	zlog.Logger.Info().Msg("shutdown server properly")
}
