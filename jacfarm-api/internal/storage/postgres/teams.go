package postgres

import (
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"context"

	"github.com/jackc/pgconn"
)

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
				return 0, storage.ErrTeamAlreadyExists
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetShortTeams(ctx context.Context) ([]*models.ShortTeam, error) {
	rows, err := s.db.Query(ctx, "SELECT id, ip FROM teams")
	if err != nil {
		return nil, err
	}

	teams := make([]*models.ShortTeam, 0)
	for rows.Next() {
		team := new(models.ShortTeam)
		err = rows.Scan(&team.ID, &team.IP)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (s *Storage) GetTeams(ctx context.Context) ([]*models.Team, error) {
	rows, err := s.db.Query(ctx, "SELECT id, name, ip FROM teams")
	if err != nil {
		return nil, err
	}

	teams := make([]*models.Team, 0)
	for rows.Next() {
		team := new(models.Team)
		err = rows.Scan(&team.ID, &team.Name, &team.IP)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}
