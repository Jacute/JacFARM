package main

import (
	"config_loader/internal/config"
	"config_loader/internal/service"
	"config_loader/internal/storage/postgres"
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
	svc.LoadConfigIntoDB(appCtx, cfg)
}
