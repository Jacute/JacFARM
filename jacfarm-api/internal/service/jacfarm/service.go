package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/pkg/rabbitmq_dto"
	"context"
	"log/slog"
)

type storage interface {
	GetFlags(ctx context.Context, filter *dto.GetFlagsFilter) ([]*models.FlagEnrich, error)
}

type queue interface {
	PublishFlag(flag *rabbitmq_dto.Flag) error
}

type Service struct {
	log *slog.Logger
	db  storage
	que queue
}

func New(log *slog.Logger, db storage, que queue) *Service {
	return &Service{
		log: log,
		db:  db,
		que: que,
	}
}
