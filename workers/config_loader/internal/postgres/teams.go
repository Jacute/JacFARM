package postgres

import (
	"config_loader/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrTeamAlreadyExists = errors.New("team already exists")

func (s *Storage) AddTeam(ctx context.Context, team *models.Team) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `INSERT INTO teams (name, ip) 
		 VALUES ($1, $2) 
		 RETURNING id`,
		team.Name, team.IP,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, ErrTeamAlreadyExists
			}
		}
		return 0, err
	}
	return id, nil
}
