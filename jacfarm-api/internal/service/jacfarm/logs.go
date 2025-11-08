package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

func (s *Service) ListLogs(ctx context.Context, filter *dto.ListLogsFilter) ([]*models.Log, int, error) {
	const op = "service.jacfarm.ListLogs"
	log := s.log.With(slog.String("op", op))

	logs, count, err := s.db.ListLogs(ctx, filter)
	if err != nil {
		log.Error("error listing logs", prettylogger.Err(err))
		return nil, 0, err
	}
	log.Debug("successfully list logs")

	return logs, count, nil
}

func (s *Service) ListLogLevel(ctx context.Context) ([]*models.LogLevel, int, error) {
	const op = "service.jacfarm.ListLogLevel"
	log := s.log.With(slog.String("op", op))

	logLevels, count, err := s.db.ListLogLevel(ctx)
	if err != nil {
		log.Error("error listing log levels", prettylogger.Err(err))
		return nil, 0, err
	}
	log.Debug("successfully list log levels")

	return logLevels, count, nil
}

func (s *Service) ListModules(ctx context.Context) ([]*models.Module, int, error) {
	const op = "service.jacfarm.ListModules"
	log := s.log.With(slog.String("op", op))

	modules, count, err := s.db.ListModules(ctx)
	if err != nil {
		log.Error("error listing modules", prettylogger.Err(err))
		return nil, 0, err
	}
	log.Debug("successfully list modules")

	return modules, count, nil
}
