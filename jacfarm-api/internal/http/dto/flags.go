package dto

import (
	"JacFARM/internal/models"
	"fmt"
	"strconv"
)

type ListFlagsFilter struct {
	Limit     uint64
	Page      uint64
	ExploitID string
	TeamID    int64
	StatusID  int64
}

type GetFlagsResponse struct {
	*Response
	Flags []*models.FlagEnrich `json:"flags"`
	Count int                  `json:"count"`
}

type PutFlagRequest struct {
	Flag string `json:"flag"`
}

type GetStatusesResponse struct {
	*Response
	Statuses []*models.Status `json:"statuses"`
}

func MapQueryToGetFlagsFilter(queries map[string]string) (*ListFlagsFilter, error) {
	exploitID := queries["exploit_id"]

	var teamID int
	teamIDStr, ok := queries["team_id"]
	if ok && teamIDStr != "" {
		var err error
		teamID, err = strconv.Atoi(teamIDStr)
		if err != nil {
			return nil, fmt.Errorf("team_id should be number")
		}
	}

	var statusID int
	statusIDStr, ok := queries["status_id"]
	if ok && statusIDStr != "" {
		var err error
		statusID, err = strconv.Atoi(statusIDStr)
		if err != nil {
			return nil, fmt.Errorf("status_id should be number")
		}
	}

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

	return &ListFlagsFilter{
		ExploitID: exploitID,
		TeamID:    int64(teamID),
		Limit:     uint64(limit),
		Page:      uint64(page),
		StatusID:  int64(statusID),
	}, nil
}
