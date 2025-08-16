package main

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"JacFARM/internal/http/server"
	"JacFARM/internal/rabbitmq"
	"JacFARM/internal/service/jacfarm"
	"JacFARM/internal/storage/postgres"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	appCtx := context.Background()
	cfg := config.MustParseConfig()

	var log *slog.Logger
	if cfg.Env == "local" {
		log = slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else if cfg.Env == "prod" {
		log = slog.New(prettylogger.NewJsonHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	} else {
		panic("invalid env parameter. should be prod|local")
	}

	// init db & rabbitmq
	db, err := postgres.New(appCtx, cfg.DB)
	if err != nil {
		panic("error connecting to db: " + err.Error())
	}
	log.Info("database connection established")
	rabbitmq := rabbitmq.New(cfg.Rabbit)

	// init farm main service
	farm := jacfarm.New(log, db, rabbitmq, cfg.ExploitDir)
	h := handlers.New(farm)
	httpServer := server.New(log, cfg.HTTP, cfg.ApiKey, h)

	go httpServer.Start()
	log.Info("jacfarm-api service started", slog.String("env", cfg.Env))

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, syscall.SIGTERM)
	<-sgn
	log.Info("shutting down JacFARM service")
	db.Stop()
	if err := rabbitmq.Close(); err != nil {
		log.Error("error closing RabbitMQ connection", prettylogger.Err(err))
	}
}
