package postgres

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
)

func (s *Storage) PutFlag(ctx context.Context, flag *models.Flag) (int64, error) {
	var lastInsertedId int64
	err := s.db.QueryRow(ctx, `
		INSERT INTO flags (value, status_id, exploit_id, get_from, message_from_server, created_at)
		VALUES ($1, (SELECT id FROM statuses WHERE name = $2), $3, $4, $5, $6)
		RETURNING id
	`, flag.Value, flag.Status, flag.ExploitID, flag.GetFrom, flag.MessageFromServer, flag.CreatedAt).Scan(&lastInsertedId)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, storage.ErrFlagAlreadyExists
			}
		}
		return 0, err
	}

	return lastInsertedId, nil
}

func (s *Storage) GetFlags(ctx context.Context, filter *dto.GetFlagsFilter) ([]*models.FlagEnrich, error) {
	builder := sq.Select("f.id", "f.value", "f.message_from_server", "f.created_at", "s.name",
		"e.id", "e.name", "e.type", "e.is_local", "e.executable_path", "e.requirements_path",
		"t.id", "t.name", "t.ip",
	).From("flags f").
		Join("exploits e ON e.id = f.exploit_id").
		Join("teams t ON t.id = f.get_from").
		Join("statuses s ON s.id = f.status_id").PlaceholderFormat(sq.Dollar)

	// apply filters
	if filter.ExploitID != "" {
		builder = builder.Where(sq.Eq{"e.id": filter.ExploitID})
	}
	if filter.TeamID != 0 {
		builder = builder.Where(sq.Eq{"t.id": filter.TeamID})
	}
	if filter.Limit > 0 {
		builder = builder.Limit(filter.Limit)
		if filter.Page > 0 {
			builder = builder.Offset(filter.Limit * (filter.Page - 1))
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error run query: %w", err)
	}

	flags := make([]*models.FlagEnrich, 0)
	for rows.Next() {
		flag := &models.FlagEnrich{
			Exploit: &models.Exploit{},
			GetFrom: &models.Team{},
		}
		err := rows.Scan(
			&flag.ID, &flag.Value, &flag.MessageFromServer, &flag.CreatedAt, &flag.Status,
			&flag.Exploit.ID, &flag.Exploit.Name, &flag.Exploit.Type, &flag.Exploit.IsLocal, &flag.Exploit.ExecutablePath, &flag.Exploit.RequirementsPath,
			&flag.GetFrom.ID, &flag.GetFrom.Name, &flag.GetFrom.IP,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan query: %w", err)
		}
		flags = append(flags, flag)
	}

	return flags, nil
}
