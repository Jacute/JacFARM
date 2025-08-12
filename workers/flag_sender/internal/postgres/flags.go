package postgres

import (
	"context"
	"errors"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	ErrFlagAlreadyExists = errors.New("flag already exists")
	ErrFlagNotUpdated    = errors.New("flag not updated")
)

func (s *Storage) PutFlag(ctx context.Context, flag *models.Flag) (int64, error) {
	var lastInsertedId int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO flags (value, status_id, exploit_id, get_from, message_from_server, created_at)
		VALUES ($1, (SELECT id FROM statuses WHERE name = $2), $3, $4, $5, $6)
		RETURNING id
	`, flag.Value, flag.Status, flag.ExploitID, flag.GetFrom, flag.MessageFromServer, flag.CreatedAt).Scan(&lastInsertedId)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, ErrFlagAlreadyExists
			}
		}
		return 0, err
	}

	return lastInsertedId, nil
}

func (s *Storage) GetFlagValuesByStatus(ctx context.Context, status models.FlagStatus) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT value
	FROM flags
	WHERE status_id = (SELECT id FROM statuses WHERE name = $1)`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flags []string
	for rows.Next() {
		var flag string
		if err := rows.Scan(&flag); err != nil {
			return nil, err
		}
		flags = append(flags, flag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return flags, nil
}

func (s *Storage) UpdateFlagByResult(ctx context.Context, flag string, result *plugins.FlagResult) error {
	res, err := s.db.Exec(ctx, `UPDATE flags
	SET status_id = (SELECT id FROM statuses WHERE name = $1),
	message_from_server = $2
	WHERE value = $3`, string(result.Status), result.Msg, flag)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return ErrFlagNotUpdated
	}
	return err
}

func (s *Storage) UpdateStatusForOldFlags(ctx context.Context, flagTTL time.Duration) (int64, error) {
	threshold := time.Now().Add(-flagTTL).UTC()

	res, err := s.db.Exec(ctx, `UPDATE flags
	SET status_id = (SELECT id FROM statuses WHERE name = $1)
	WHERE status_id = (SELECT id FROM statuses WHERE name = $2) AND
	created_at <= $3`,
		models.FlagStatusOld, models.FlagStatusPending, threshold)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
