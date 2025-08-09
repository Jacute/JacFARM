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

const dbPath = "./database.db"

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
	db, err := sqlite.New(dbPath)
	if err != nil {
		panic("error connecting to db: " + err.Error())
	}
	db.ApplyMigrations(appCtx, dbPath, cfg.DB.MigrationsPath)
	log.Info("database connection established")
	rabbitmq := rabbitmq.New(cfg.Rabbit)

	// init farm main service
	farm := jacfarm.New(log, db)
	farm.LoadConfigIntoDB(appCtx, cfg)

	// init all workers
	exploitRunner, err := exploit_runner.New(log, db, rabbitmq, cfg.ExploitRunner.ExploitDirectory)
	if err != nil {
		panic(err)
	}
	flagSaver := flag_saver.New(log, rabbitmq, db)
	flagSender, err := flag_sender.New(log, db, cfg.FlagSender.PluginDir)
	if err != nil {
		panic(err)
	}

	// run all workers
	go exploitRunner.Start(appCtx)
	go func() {
		err := flagSaver.Start()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := flagSender.Start()
		if err != nil {
			panic(err)
		}
	}()

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
