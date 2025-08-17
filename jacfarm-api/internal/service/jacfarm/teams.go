package jacfarm

import (
	"JacFARM/internal/models"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

func (s *Service) ListShortTeams(ctx context.Context) ([]*models.ShortTeam, error) {
	const op = "service.jacfarm.ListShortTeams"
	log := s.log.With(slog.String("op", op))

	teams, err := s.db.GetShortTeams(ctx)
	if err != nil {
		log.Error("error listing short teams", prettylogger.Err(err))
		return nil, err
	}

	return teams, nil
}
