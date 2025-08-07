package rabbitmq

import "JacFARM/internal/models"

type Flag struct {
	Value   string          `json:"value"`
	Exploit *models.Exploit `json:"exploit"`
	GetFrom *models.Team    `json:"victim_team"`
	IsLocal bool            `json:"is_local"`
}
