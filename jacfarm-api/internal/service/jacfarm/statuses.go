package jacfarm

import (
	"JacFARM/internal/models"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

func (s *Service) GetStatuses(ctx context.Context) ([]*models.Status, error) {
	const op = "service.jacfarm.GetStatuses"
	log := s.log.With(slog.String("op", op))

	statuses, err := s.db.GetStatuses(ctx)
	if err != nil {
		log.Error("error getting statuses", prettylogger.Err(err))
		return nil, err
	}
	log.Debug("got statuses successfully")

	return statuses, nil
}
