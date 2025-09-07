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

func (s *Storage) PutFlags(ctx context.Context, flags []*models.Flag) ([]int64, error) {
	insertBuilder := sq.Insert("flags").
		Columns("value", "status_id", "message_from_server", "created_at", "exploit_id", "get_from").
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (value) DO NOTHING RETURNING id")

	for _, flag := range flags {
		values := []interface{}{
			flag.Value,
			sq.Expr("(SELECT id FROM statuses WHERE name = ?)", flag.Status),
			flag.MessageFromServer,
			flag.CreatedAt,
			nil, // exploit_id
			nil, // get_from
		}

		if flag.ExploitID != nil {
			values[4] = *flag.ExploitID
		}
		if flag.GetFrom != nil {
			values[5] = *flag.GetFrom
		}

		insertBuilder = insertBuilder.Values(values...)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building batch insert query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing batch insert: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
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
