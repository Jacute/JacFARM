package sqlite

import (
	"JacFARM/internal/storage"
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) AddConfigParameter(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO config (key, value)
						VALUES ($1, $2)`, key, value)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return storage.ErrConfigParamAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (s *Storage) GetConfigParameter(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM config WHERE key = $1`, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
