package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

func (s *Service) GetFlags(ctx context.Context, filter *dto.GetFlagsFilter) ([]*models.FlagEnrich, error) {
	const op = "service.jacfarm.GetFlags"
	log := s.log.With(slog.String("op", op))

	flags, err := s.db.GetFlags(ctx, filter)
	if err != nil {
		log.Error("error getting flags", prettylogger.Err(err))
		return nil, err
	}

	log.Info("got flags successfully")

	return flags, nil
}

func (s *Service) PutFlag(ctx context.Context, flag string) error {
	panic("not implemented")
}
