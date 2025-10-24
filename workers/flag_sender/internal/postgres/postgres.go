package postgres

import (
	"context"
	"flag_sender/internal/config"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.db
}

func (s *Storage) Stop() {
	s.db.Close()
}

func New(ctx context.Context, config *config.DBConfig) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.DBName)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		panic("Failed to create connection pool: " + err.Error())
	}

	err = db.Ping(ctx)
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	return &Storage{db}
}
