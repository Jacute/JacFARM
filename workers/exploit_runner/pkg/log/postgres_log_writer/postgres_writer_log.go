package postgreslogwriter

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	builder := sq.Insert("audit.logs").
		Columns("value", "created_at", "log_level_id")
	values := []any{value, createdAt, sq.Expr("(SELECT id FROM audit.log_levels WHERE name = ?)", level)}
	if module != "" {
		builder = builder.Columns("module_id")
		values = append(values, sq.Expr("(SELECT id FROM audit.modules WHERE name = ?)", module))
	}
	if op != "" {
		builder = builder.Columns("operation")
		values = append(values, op)
	}
	if exploitId != "" {
		builder = builder.Columns("exploit_id")
		values = append(values, exploitId)
	}
	if attrs != nil {
		builder = builder.Columns("attrs")
		values = append(values, attrs)
	}
	builder = builder.Values(values...).PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return ErrWriteLog
	}

	cmd, err := plg.db.Exec(ctx, query, args...)
	if err != nil {
		return ErrWriteLog
	}
	if cmd.RowsAffected() == 0 {
		return ErrWriteLog
	}

	return nil
}
