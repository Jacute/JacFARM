package main

import (
	"config_loader/internal/config"
	"config_loader/internal/postgres"
	"config_loader/internal/service"
	"context"
	"log/slog"
	"os"

	"github.com/jacute/prettylogger"
)

func main() {
	appCtx := context.Background()
	log := slog.New(prettylogger.NewJsonHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := config.MustParseConfig()

	db := postgres.New(appCtx, cfg.DB)
	defer db.Stop()

	svc := service.New(log, db)
	err := svc.LoadConfigIntoDB(appCtx, cfg)
	if err != nil {
		panic(err)
	}
	log.Info("config loaded into db successfully")
}
