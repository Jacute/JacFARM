package postgres

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetFlags(ctx context.Context, filter *dto.ListFlagsFilter) ([]*models.FlagEnrich, int, error) {
	baseBuilder := sq.Select("f.id", "f.value", "f.message_from_server", "f.created_at",
		"s.name", "COALESCE(e.name, '')", "COALESCE(t.ip::text, '')").
		From("flags f").
		LeftJoin("exploits e ON e.id = f.exploit_id").
		LeftJoin("teams t ON t.id = f.get_from").
		Join("statuses s ON s.id = f.status_id").
		PlaceholderFormat(sq.Dollar)

	// apply filters
	if filter.ExploitID != "" {
		baseBuilder = baseBuilder.Where(sq.Eq{"e.id": filter.ExploitID})
	}
	if filter.TeamID != 0 {
		baseBuilder = baseBuilder.Where(sq.Eq{"t.id": filter.TeamID})
	}

	countBuilder := sq.Select("COUNT(*)").
		From("flags f").
		LeftJoin("exploits e ON e.id = f.exploit_id").
		LeftJoin("teams t ON t.id = f.get_from").
		Join("statuses s ON s.id = f.status_id").
		PlaceholderFormat(sq.Dollar)

	if filter.ExploitID != "" {
		countBuilder = countBuilder.Where(sq.Eq{"e.id": filter.ExploitID})
	}
	if filter.TeamID != 0 {
		countBuilder = countBuilder.Where(sq.Eq{"t.id": filter.TeamID})
	}

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("error building count query: %w", err)
	}

	var total int
	err = s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error running count query: %w", err)
	}

	if filter.Limit > 0 {
		baseBuilder = baseBuilder.Limit(filter.Limit)
		if filter.Page > 0 {
			baseBuilder = baseBuilder.Offset(filter.Limit * (filter.Page - 1))
		}
	}

	query, args, err := baseBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("error building query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error run query: %w", err)
	}
	defer rows.Close()

	flags := make([]*models.FlagEnrich, 0)
	for rows.Next() {
		flag := &models.FlagEnrich{}
		err := rows.Scan(
			&flag.ID, &flag.Value, &flag.MessageFromServer, &flag.CreatedAt,
			&flag.Status, &flag.ExploitName, &flag.VictimIP,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scan query: %w", err)
		}
		flags = append(flags, flag)
	}

	return flags, total, nil
}
