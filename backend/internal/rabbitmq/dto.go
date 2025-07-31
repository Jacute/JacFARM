package rabbitmq

import "JacFARM/internal/models"

type Flag struct {
	Value       string       `json:"value"`
	ExploitName string       `json:"exploit_name"`
	GetFrom     *models.Team `json:"victim_team"`
}
