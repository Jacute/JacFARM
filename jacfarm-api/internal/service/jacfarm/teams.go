package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	storage_errors "JacFARM/internal/storage"
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

	log.Info("short teams successfully listed")

	return teams, nil
}

func (s *Service) ListTeams(ctx context.Context, filter *dto.ListTeamsFilter) ([]*models.Team, int, error) {
	const op = "service.jacfarm.ListTeams"
	log := s.log.With(slog.String("op", op))

	teams, count, err := s.db.GetTeams(ctx, filter)
	if err != nil {
		log.Error("error listing teams", prettylogger.Err(err))
		return nil, 0, err
	}

	slog.Info("teams successfully listed", slog.Int("count", count))

	return teams, count, nil
}

func (s *Service) AddTeam(ctx context.Context, team *models.Team) (int64, error) {
	const op = "service.jacfarm.AddTeam"
	log := s.log.With(slog.String("op", op))

	id, err := s.db.AddTeam(ctx, team)
	if err != nil {
		if err == storage_errors.ErrTeamAlreadyExists {
			return 0, err
		}
		log.Error("error adding team", prettylogger.Err(err))
		return 0, err
	}

	log.Debug("team successfully added", slog.Int64("id", id))

	return id, nil
}

func (s *Service) DeleteTeam(ctx context.Context, id int64) error {
	const op = "service.jacfarm.DeleteTeam"
	log := s.log.With(slog.String("op", op), slog.Int64("id", id))

	err := s.db.DeleteTeam(ctx, id)
	if err != nil {
		if err == storage_errors.ErrTeamNotFound {
			return err
		}
		log.Error("error deleting team", prettylogger.Err(err))
		return err
	}

	log.Debug("team successfully deleted")

	return nil
}
