package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var ErrConfigParamAlreadyExists = errors.New("config param already exists")

func (s *Storage) AddConfigParameter(ctx context.Context, key, value string) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `INSERT INTO config (key, value)
		VALUES ($1, $2) RETURNING id`, key, value).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, ErrConfigParamAlreadyExists
			}
		}
		return 0, err
	}

	return id, nil
}
