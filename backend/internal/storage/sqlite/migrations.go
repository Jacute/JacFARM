package sqlite

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (s *Storage) ApplyMigrations(ctx context.Context, dbPath string, migrationsPath string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		"sqlite3://"+dbPath,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		panic(err)
	}
}
