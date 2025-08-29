package postgres

import (
	"context"
	"errors"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var (
	ErrFlagAlreadyExists = errors.New("flag already exists")
	ErrFlagNotUpdated    = errors.New("flag not updated")
)

func (s *Storage) PutFlag(ctx context.Context, flag *models.Flag) (int64, error) {
	insertBuilder := sq.Insert("flags").PlaceholderFormat(sq.Dollar)
	insertMap := map[string]interface{}{
		"value":               flag.Value,
		"status_id":           sq.Expr("(SELECT id FROM statuses WHERE name = ?)", flag.Status),
		"message_from_server": flag.MessageFromServer,
		"created_at":          flag.CreatedAt,
	}
	if flag.ExploitID != nil {
		insertMap["exploit_id"] = *flag.ExploitID
	}
	if flag.GetFrom != nil {
		insertMap["get_from"] = *flag.GetFrom
	}
	insertBuilder = insertBuilder.SetMap(insertMap).Suffix("ON CONFLICT (value) DO NOTHING RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building query: %s", err.Error())
	}

	var lastInsertedId int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&lastInsertedId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrFlagAlreadyExists
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
