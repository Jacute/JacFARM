package postgres

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetFlags(ctx context.Context, filter *dto.ListFlagsFilter) ([]*models.FlagEnrich, error) {
	builder := sq.Select("f.id", "f.value", "f.message_from_server", "f.created_at", "s.name", "COALESCE(e.name, '')", "COALESCE(t.ip::text, '')").From("flags f").
		LeftJoin("exploits e ON e.id = f.exploit_id").
		LeftJoin("teams t ON t.id = f.get_from").
		Join("statuses s ON s.id = f.status_id").PlaceholderFormat(sq.Dollar)

	// apply filters
	if filter.ExploitID != "" {
		builder = builder.Where(sq.Eq{"e.id": filter.ExploitID})
	}
	if filter.TeamID != 0 {
		builder = builder.Where(sq.Eq{"t.id": filter.TeamID})
	}
	if filter.Limit > 0 {
		builder = builder.Limit(filter.Limit)
		if filter.Page > 0 {
			builder = builder.Offset(filter.Limit * (filter.Page - 1))
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error run query: %w", err)
	}

	flags := make([]*models.FlagEnrich, 0)
	for rows.Next() {
		flag := &models.FlagEnrich{}
		err := rows.Scan(
			&flag.ID, &flag.Value, &flag.MessageFromServer, &flag.CreatedAt,
			&flag.Status, &flag.ExploitName, &flag.VictimIP,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan query: %w", err)
		}
		flags = append(flags, flag)
	}

	return flags, nil
}
