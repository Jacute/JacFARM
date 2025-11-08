package postgres

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) ListLogs(ctx context.Context, filter *dto.ListLogsFilter) ([]*models.Log, int, error) {
	builder := sq.Select(
		"l.id",
		"m.name",
		"l.operation",
		"e.name",
		"l.value",
		"l.attrs",
		"l.created_at",
		"ll.name",
	).From("audit.logs l").
		LeftJoin("audit.modules m on l.module_id = m.id").
		LeftJoin("exploits e on e.id = l.exploit_id").
		LeftJoin("audit.log_levels ll on l.log_level_id = ll.id").
		OrderBy("l.created_at DESC").
		PlaceholderFormat(sq.Dollar)
	countBuilder := sq.Select("COUNT(*)").From("audit.logs l").PlaceholderFormat(sq.Dollar)

	if filter.Limit > 0 {
		builder = builder.Limit(uint64(filter.Limit))
		if filter.Page > 0 {
			builder = builder.Offset(uint64(filter.Limit * (filter.Page - 1)))
		}
	}

	if filter.ModuleId != 0 {
		builder = builder.Where(sq.Eq{"l.module_id": filter.ModuleId})
		countBuilder = countBuilder.Where(sq.Eq{"l.module_id": filter.ModuleId})
	}
	if filter.ExploitId != "" {
		builder = builder.Where(sq.Eq{"l.exploit_id": filter.ExploitId})
		countBuilder = countBuilder.Where(sq.Eq{"l.exploit_id": filter.ExploitId})
	}
	if filter.LogLevelId != 0 {
		builder = builder.Where(sq.Eq{"l.log_level_id": filter.LogLevelId})
		countBuilder = countBuilder.Where(sq.Eq{"l.log_level_id": filter.LogLevelId})
	}

	builder.ToSql()
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("error building query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		if err := rows.Scan(&log.Id, &log.Module, &log.Operation, &log.Exploit, &log.Value, &log.Attrs, &log.CreatedAt, &log.Level); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("error building count query: %w", err)
	}

	var count int
	if err := s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&count); err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}

func (s *Storage) ListModules(ctx context.Context) ([]*models.Module, int, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name FROM audit.modules`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	modules := make([]*models.Module, 0)
	for rows.Next() {
		var module models.Module
		if err := rows.Scan(&module.Id, &module.Name); err != nil {
			return nil, 0, err
		}
		modules = append(modules, &module)
	}

	return modules, len(modules), nil
}

func (s *Storage) ListLogLevel(ctx context.Context) ([]*models.LogLevel, int, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name FROM audit.log_levels`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logLevels := make([]*models.LogLevel, 0)
	for rows.Next() {
		var level models.LogLevel
		if err := rows.Scan(&level.Id, &level.Name); err != nil {
			return nil, 0, err
		}
		logLevels = append(logLevels, &level)
	}

	return logLevels, len(logLevels), nil
}
