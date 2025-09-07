package main

import (
	"context"
	"flag_sender/internal/config"
	"flag_sender/internal/postgres"
	"flag_sender/internal/rabbitmq"
	"flag_sender/internal/services/flag_sender"
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
	db := postgres.New(appCtx, cfg.DB)
	log.Info("database connection established")
	rabbitmq := rabbitmq.New(cfg.Rabbit)

	// init services
	flagSender, err := flag_sender.New(log, cfg.PluginDir, rabbitmq, db)
	if err != nil {
		panic(err)
	}
	go func() {
		err := flagSender.Start()
		if err != nil {
			panic(err)
		}
	}()
	log.Info("flag saver & flag sender services started", slog.String("env", cfg.Env))

	// graceful shutdown
	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, syscall.SIGTERM)
	<-sgn
	log.Info("shutting down service")
	flagSender.Stop()
	db.Stop()
	if err := rabbitmq.Close(); err != nil {
		log.Error("error closing RabbitMQ connection", prettylogger.Err(err))
	}
}
