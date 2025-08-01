package main

import (
	"JacFARM/internal/config"
	"JacFARM/internal/rabbitmq"
	jacfarm "JacFARM/internal/service/JacFARM"
	"JacFARM/internal/service/exploit_runner"
	"JacFARM/internal/service/flag_saver"
	"JacFARM/internal/service/flag_sender"
	"JacFARM/internal/storage/sqlite"
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

	db, err := sqlite.New()
	if err != nil {
		panic("error connecting to db: " + err.Error())
	}
	db.ApplyMigrations(appCtx, cfg.DB.MigrationsPath)
	log.Info("database connection established")

	rabbitmq := rabbitmq.New(cfg.Rabbit)

	farm := jacfarm.New(log, db)
	farm.LoadConfigIntoDB(appCtx, cfg)

	exploitRunner := exploit_runner.New(log, db, rabbitmq, cfg.ExploitRunner.ExploitDirectory)
	flagSaver := flag_saver.New(log, rabbitmq, db)
	flagSender := flag_sender.New(log, db)
	go exploitRunner.Start(appCtx)
	go flagSaver.Start()
	go flagSender.Start()

	log.Info("JacFARM service started", slog.String("env", cfg.Env))

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, syscall.SIGTERM)
	<-sgn
	log.Info("shutting down JacFARM service")
	exploitRunner.Stop()
	flagSaver.Stop()
	flagSender.Stop()
	if err := db.Close(); err != nil {
		log.Error("error closing database connection", prettylogger.Err(err))
	}
	if err := rabbitmq.Close(); err != nil {
		log.Error("error closing RabbitMQ connection", prettylogger.Err(err))
	}
}
