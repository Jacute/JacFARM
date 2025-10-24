package main

import (
	"context"
	"flag_sender/internal/config"
	"flag_sender/internal/postgres"
	"flag_sender/internal/rabbitmq"
	"flag_sender/internal/services/flag_sender"
	postgreslogwriter "flag_sender/pkg/log/postgres_log_writer"
	postgresslog "flag_sender/pkg/log/postgres_slog"
	slogmultihandler "flag_sender/pkg/log/slog_multihandler"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

const moduleName = "flag_sender"

func main() {
	appCtx := context.Background()
	cfg := config.MustParseConfig()

	// init db
	db := postgres.New(appCtx, cfg.DB)

	// init logger
	var options *slog.HandlerOptions

	handlers := make([]slog.Handler, 0, 2)
	if cfg.Env == "local" {
		options = &slog.HandlerOptions{Level: slog.LevelDebug}
		handlers = append(handlers, prettylogger.NewColoredHandler(os.Stdout, options))
	} else if cfg.Env == "prod" {
		options = &slog.HandlerOptions{Level: slog.LevelInfo}
		handlers = append(handlers, prettylogger.NewJsonHandler(os.Stdout, options))
	} else {
		panic("invalid env parameter. should be prod|local")
	}
	postgresHandler := postgresslog.NewHandler(moduleName, postgreslogwriter.New(db.GetPool()), options)
	handlers = append(handlers, postgresHandler)
	log := slog.New(slogmultihandler.New(handlers))

	// init rabbitmq
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
