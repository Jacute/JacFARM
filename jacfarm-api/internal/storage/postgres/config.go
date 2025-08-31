package postgres

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetConfig(ctx context.Context, filter *dto.GetConfigFilter) ([]*models.Config, int, error) {
	builder := sq.Select("id", "key", "value").
		From("config").
		PlaceholderFormat(sq.Dollar)
	if filter.Limit > 0 {
		builder = builder.Limit(filter.Limit)
		if filter.Page > 0 {
			builder = builder.Offset(filter.Limit * (filter.Page - 1))
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("error building query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var config []*models.Config
	for rows.Next() {
		var configRow models.Config
		if err := rows.Scan(&configRow.ID, &configRow.Name, &configRow.Value); err != nil {
			return nil, 0, err
		}
		config = append(config, &configRow)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var count int
	if err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM config").Scan(&count); err != nil {
		return nil, 0, err
	}

	return config, count, nil
}

func (s *Storage) UpdateConfigRow(ctx context.Context, id int64, value string) error {
	cmd, err := s.db.Exec(ctx, "UPDATE config SET value = $1 WHERE id = $2", value, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return storage.ErrConfigParamNotFound
	}
	return nil
}
