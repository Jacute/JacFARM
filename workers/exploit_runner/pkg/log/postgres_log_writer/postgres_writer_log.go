package postgreslogwriter

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrWriteLog = errors.New("error write log into logs table in postgres")
)

type PostgresLogWriter struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresLogWriter {
	return &PostgresLogWriter{db}
}

func (plg *PostgresLogWriter) WriteLog(ctx context.Context, module, op, level, value, exploitId string, attrs map[string]any, createdAt time.Time) error {
	cmd, err := plg.db.Exec(ctx, `INSERT INTO audit.logs
	(module_id, operation, value, attrs, created_at, log_level_id)
	VALUES (
	(SELECT id FROM audit.modules WHERE name = $1),
	$2, $3, $4, $5,
	(SELECT id FROM audit.log_levels WHERE name = $6)
	)`,
		module, op, value, attrs, createdAt, level)

	if err != nil {
		return ErrWriteLog
	}
	if cmd.RowsAffected() == 0 {
		return ErrWriteLog
	}

	return nil
}
