package dto

import (
	"JacFARM/internal/models"
	"fmt"
	"strconv"
)

type ListShortTeamsResponse struct {
	*Response
	Teams []*models.ShortTeam `json:"teams"`
}

type ListTeamsFilter struct {
	Page  uint64
	Limit uint64
	Query string
}

type ListTeamsResponse struct {
	*Response
	Teams []*models.Team `json:"teams"`
	Count int            `json:"count"`
}

type AddTeamRequest struct {
	Name string `json:"name"`
	IP   string `json:"ip" validate:"required,ip"`
}

type AddTeamResponse struct {
	*Response
	ID int64 `json:"id"`
}

func MapQueryToListTeamsFilter(queries map[string]string) (*ListTeamsFilter, error) {
	var limit int
	limitStr, ok := queries["limit"]
	if ok && limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("limit should be number")
		}
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit should be positive number")
	}

	var page int
	pageStr, ok := queries["page"]
	if ok && pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("page should be number")
		}
	}
	if page < 0 {
		return nil, fmt.Errorf("page should be positive number")
	}

	return &ListTeamsFilter{
		Limit: uint64(limit),
		Page:  uint64(page),
		Query: queries["query"],
	}, nil
}
