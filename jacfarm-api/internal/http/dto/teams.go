package dto

import "JacFARM/internal/models"

type ListShortTeamsResponse struct {
	*Response
	Teams []*models.ShortTeam `json:"teams"`
}
