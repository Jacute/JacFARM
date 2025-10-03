package postgres

import (
	"JacFARM/internal/models"
	"context"
)

func (s *Storage) GetStatuses(ctx context.Context) ([]*models.Status, error) {
	rows, err := s.db.Query(ctx, "SELECT id, name FROM statuses")
	if err != nil {
		return nil, err
	}

	statuses := make([]*models.Status, 0)
	for rows.Next() {
		status := new(models.Status)
		err = rows.Scan(&status.ID, &status.Name)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}
