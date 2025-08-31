package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

func (s *Service) GetConfig(ctx context.Context, filter *dto.GetConfigFilter) ([]*models.Config, int, error) {
	const op = "service.jacfarm.GetConfig"
	log := s.log.With(slog.String("op", op))

	config, count, err := s.db.GetConfig(ctx, filter)
	if err != nil {
		log.Error("error getting config", prettylogger.Err(err))
		return nil, 0, err
	}
	log.Debug("got config successfully")

	return config, count, nil
}

func (s *Service) UpdateConfig(ctx context.Context, id int64, value string) error {
	const op = "service.jacfarm.PutConfig"
	log := s.log.With(slog.String("op", op))

	err := s.db.UpdateConfigRow(ctx, id, value)
	if err != nil {
		log.Error("error updating config", prettylogger.Err(err))
		return err
	}
	log.Debug("config successfully updated")

	return nil
}
