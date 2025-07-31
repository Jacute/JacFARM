package sqlite

import "JacFARM/internal/models"

func (s *Storage) PutFlag(flag *models.Flag) error {
	_, err := s.db.Exec(`INSERT INTO flags (value, status_id, exploit_id, get_from, message_from_server)
	VALUES ($1, (SELECT id FROM statuses WHERE name = $2), $3, $4, $5)`, flag.Value, flag.Status, flag.Exploit.ID, flag.GetFrom.ID, flag.MessageFromServer)
	return err
}
