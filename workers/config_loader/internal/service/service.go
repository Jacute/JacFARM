package service

import (
	"config_loader/internal/models"
	"context"
	"log/slog"
)

type storage interface {
	AddTeam(ctx context.Context, team *models.Team) (int64, error)
	AddConfigParameter(ctx context.Context, key, value string) (int64, error)
}

type Service struct {
	db  storage
	log *slog.Logger
}

func New(log *slog.Logger, db storage) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}
