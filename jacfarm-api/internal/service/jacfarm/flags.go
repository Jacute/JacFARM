package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/pkg/rabbitmq_dto"
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
	const op = "service.jacfarm.PutFlag"
	log := s.log.With(slog.String("op", op), slog.String("flag", flag))

	err := s.que.PublishFlag(&rabbitmq_dto.Flag{
		Value:      flag,
		ExploitID:  "",
		TeamID:     0,
		SourceType: rabbitmq_dto.ManualSendingSourceType,
	})
	if err != nil {
		log.Error("error sending flag to queue", prettylogger.Err(err))
		return err
	}
	log.Info("flag successfully send to queue")

	return nil
}
