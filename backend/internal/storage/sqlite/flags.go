package sqlite

import (
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"context"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) PutFlag(ctx context.Context, flag *models.Flag) (int64, error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO flags
	(value, status_id, exploit_id, get_from, message_from_server, created_at)
	VALUES ($1, (SELECT id FROM statuses WHERE name = $2), $3, $4, $5, $6)`, flag.Value, flag.Status, flag.ExploitID, flag.GetFrom, flag.MessageFromServer, flag.CreatedAt)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return 0, storage.ErrFlagAlreadyExists
			}
		}
		return 0, err
	}
	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertedId, err
}

func (s *Storage) GetFlagsByStatus(ctx context.Context, status models.FlagStatus) ([]*models.Flag, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, value, status_id, exploit_id, get_from, message_from_server
	FROM flags
	WHERE status_id = (SELECT id FROM statuses WHERE name = $1)`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flags []*models.Flag
	for rows.Next() {
		var flag models.Flag
		if err := rows.Scan(&flag.ID, &flag.Value, &flag.Status, &flag.ExploitID, &flag.GetFrom, &flag.MessageFromServer); err != nil {
			return nil, err
		}
		flags = append(flags, &flag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return flags, nil
}

func (s *Storage) GetFlagValuesByStatus(ctx context.Context, status models.FlagStatus) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT value
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

func (s *Storage) UpdateStatusForOldFlags(ctx context.Context, flagTTL time.Duration) (int64, error) {
	threshold := time.Now().Add(-flagTTL).UTC().Unix()

	res, err := s.db.ExecContext(ctx, `UPDATE flags
	SET status_id = (SELECT id FROM statuses WHERE name = $1)
	WHERE status_id = (SELECT id FROM statuses WHERE name = $2) AND
	created_at <= $3`,
		models.FlagStatusOld, models.FlagStatusPending, threshold)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, err
}

func (s *Storage) GetFlagEnrichByValue(ctx context.Context, value string) (*models.FlagEnrich, error) {
	flag := models.FlagEnrich{
		Exploit: &models.Exploit{},
		GetFrom: &models.Team{},
	}
	err := s.db.QueryRowContext(ctx, `SELECT f.id, f.value, s.name, f.message_from_server,
	e.id, e.name, e.type, e.is_local, e.executable_path, e.requirements_path, t.id, t.name, t.ip
	FROM flags f
	JOIN statuses s ON f.status_id = s.id
	JOIN exploits e ON f.exploit_id = e.id
	JOIN teams t ON f.get_from = t.id
	WHERE f.value = $1`, value).
		Scan(&flag.ID, &flag.Value, &flag.Status, &flag.MessageFromServer,
			&flag.Exploit.ID, &flag.Exploit.Name, &flag.Exploit.Type, &flag.Exploit.IsLocal, &flag.Exploit.ExecutablePath, &flag.Exploit.RequirementsPath,
			&flag.GetFrom.ID, &flag.GetFrom.Name, &flag.GetFrom.IP)
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

func (s *Storage) UpdateFlagStatus(ctx context.Context, flag string, status models.FlagStatus) error {
	res, err := s.db.ExecContext(ctx, `UPDATE flags
	SET status_id = (SELECT id FROM statuses WHERE name = $1)
	WHERE value = $2`, string(status), flag)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return storage.ErrFlagNotUpdated
	}
	return err
}
