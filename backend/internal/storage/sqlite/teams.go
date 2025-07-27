package sqlite

import (
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"context"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) AddTeam(team *models.Team) error {
	_, err := s.db.Exec(`INSERT INTO teams (name, ip)
	VALUES ($1, $2)`, team.Name, team.IP)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return storage.ErrTeamAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (s *Storage) GetTeams(ctx context.Context) ([]*models.Team, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, ip FROM teams")
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
