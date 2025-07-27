package sqlite

import (
	"context"
)

func (s *Storage) AddConfigParameter(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx, `INSERT OR REPLACE INTO config (key, value)
						VALUES ($1, $2)`, key, value)
	if err != nil {
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
