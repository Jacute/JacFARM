package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"log/slog"
)

type storage interface {
	GetFlags(ctx context.Context, filter *dto.GetFlagsFilter) ([]*models.FlagEnrich, error)
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
