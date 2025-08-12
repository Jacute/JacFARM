package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrConfigParamAlreadyExists = errors.New("config param already exists")

func (s *Storage) AddConfigParameter(ctx context.Context, key, value string) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `INSERT INTO config (key, value)
		VALUES ($1, $2) RETURNING id`, key, value).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, ErrConfigParamAlreadyExists
			}
		}
		return 0, err
	}

	return id, nil
}
