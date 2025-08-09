package postgres

import (
	"config_loader/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var ErrTeamAlreadyExists = errors.New("team already exists")

func (s *Storage) AddTeam(ctx context.Context, team *models.Team) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `INSERT INTO teams (name, ip) 
		 VALUES ($1, $2) 
		 RETURNING id`,
		team.Name, team.IP,
	).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return 0, ErrTeamAlreadyExists
		}
	}
	return id, nil
}
