package dto

import (
	"JacFARM/internal/models"
	"fmt"
	"strconv"
)

type GetConfigResponse struct {
	*Response
	Config []*models.Config `json:"config"`
	Count  int              `json:"count"`
}

type GetConfigFilter struct {
	Limit uint64 `json:"limit"`
	Page  uint64 `json:"page"`
}

type UpdateConfigRequest struct {
	Value string `json:"value"`
}

func MapQueryToGetConfigFilter(queries map[string]string) (*GetConfigFilter, error) {
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

	return &GetConfigFilter{
		Limit: uint64(limit),
		Page:  uint64(page),
	}, nil
}
