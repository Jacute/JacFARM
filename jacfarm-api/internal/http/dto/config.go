package dto

import (
	"JacFARM/internal/models"
	"fmt"
	"strconv"
)

var (
	ErrLimitIncorrectType = fmt.Errorf("limit should be number")
	ErrPageIncorrectType  = fmt.Errorf("page should be number")
	ErrLimitNegative      = fmt.Errorf("limit should be positive number")
	ErrPageNegative       = fmt.Errorf("page should be positive number")
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
			return nil, ErrLimitIncorrectType
		}
	}
	if limit < 0 {
		return nil, ErrLimitNegative
	}

	var page int
	pageStr, ok := queries["page"]
	if ok && pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, ErrPageIncorrectType
		}
	}
	if page < 0 {
		return nil, ErrPageNegative
	}

	return &GetConfigFilter{
		Limit: uint64(limit),
		Page:  uint64(page),
	}, nil
}
