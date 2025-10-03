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

func (s *Storage) GetTeams(ctx context.Context, filter *dto.ListTeamsFilter) ([]*models.Team, int, error) {
	builder := sq.Select("id", "name", "ip").From("teams").PlaceholderFormat(sq.Dollar)
	countBuilder := sq.Select("COUNT(*)").From("teams").PlaceholderFormat(sq.Dollar)

	if filter.Limit > 0 {
		builder = builder.Limit(filter.Limit)
		if filter.Page > 0 {
			builder = builder.Offset(filter.Limit * (filter.Page - 1))
		}
	}

	if filter.Query != "" {
		builder = builder.Where(sq.Or{
			sq.ILike{"name": "%" + filter.Query + "%"},
			sq.Expr("ip::text ILIKE ?", "%"+filter.Query+"%"),
		})
		countBuilder = countBuilder.Where(sq.Or{
			sq.ILike{"name": "%" + filter.Query + "%"},
			sq.Expr("ip::text ILIKE ?", "%"+filter.Query+"%"),
		})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error query: %w", err)
	}

	teams := make([]*models.Team, 0)
	for rows.Next() {
		team := new(models.Team)
		err = rows.Scan(&team.ID, &team.Name, &team.IP)
		if err != nil {
			return nil, 0, err
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error in rows: %w", err)
	}

	query, args, err = countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build sql query: %w", err)
	}

	var count int
	err = s.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("error query: %w", err)
	}

	return teams, count, nil
}

func (s *Storage) DeleteTeam(ctx context.Context, id int64) error {
	cmd, err := s.db.Exec(ctx, "DELETE FROM teams WHERE id = $1", id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return storage.ErrTeamNotFound
	}

	return nil
}
