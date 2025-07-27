package jacfarm

import (
	"JacFARM/internal/models"
	"context"
	"log/slog"
)

type storage interface {
	PutFlag(flag *models.Flag) error
	AddTeam(team *models.Team) error
	AddConfigParameter(ctx context.Context, key, value string) error
}

type Service struct {
	log *slog.Logger
	db  storage
}

func New(log *slog.Logger, db storage) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}
