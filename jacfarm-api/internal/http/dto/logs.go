package dto

import (
	"JacFARM/internal/models"
	"errors"
	"strconv"

	"github.com/google/uuid"
)

type ListLogsFilter struct {
	Page       int
	Limit      int
	ExploitId  string
	ModuleId   int
	LogLevelId int
}

type ListLogsResponse struct {
	*Response
	Logs  []*models.Log `json:"logs"`
	Count int           `json:"count"`
}

var (
	ErrModuleIdIncorrect   = errors.New("module_id should be number")
	ErrExploitIdIncorrect  = errors.New("exploit_id should be UUID")
	ErrLogLevelIdIncorrect = errors.New("log_level_id should be number")

	ErrModuleIdNegative   = errors.New("module_id should be positive number")
	ErrLogLevelIdNegative = errors.New("log_level_id should be positive number")
)

func MapQueryToListLogsFilter(queries map[string]string) (*ListLogsFilter, error) {
	var limit, page, moduleId, logLevelId int
	var exploitId string
	limitStr, ok := queries["limit"]
	if ok && limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, ErrLimitIncorrect
		}
		if limit < 0 {
			return nil, ErrLimitNegative
		}
	}

	pageStr, ok := queries["page"]
	if ok && pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, ErrPageIncorrect
		}
		if page < 0 {
			return nil, ErrPageNegative
		}
	}

	exploitIdStr, ok := queries["exploit_id"]
	if ok && exploitIdStr != "" {
		var err error
		exploitIdUUID, err := uuid.Parse(exploitIdStr)
		if err != nil {
			return nil, ErrExploitIdIncorrect
		}
		exploitId = exploitIdUUID.String()
	}

	moduleIdStr, ok := queries["module_id"]
	if ok && moduleIdStr != "" {
		var err error
		moduleId, err = strconv.Atoi(moduleIdStr)
		if err != nil {
			return nil, ErrModuleIdIncorrect
		}
		if moduleId < 0 {
			return nil, ErrModuleIdNegative
		}
	}

	logLevelIdStr, ok := queries["log_level_id"]
	if ok && logLevelIdStr != "" {
		var err error
		logLevelId, err = strconv.Atoi(logLevelIdStr)
		if err != nil {
			return nil, ErrLogLevelIdIncorrect
		}
		if logLevelId < 0 {
			return nil, ErrLogLevelIdNegative
		}
	}

	return &ListLogsFilter{
		Page:       page,
		Limit:      limit,
		ExploitId:  exploitId,
		ModuleId:   moduleId,
		LogLevelId: logLevelId,
	}, nil
}

type ListModulesResponse struct {
	*Response
	Modules []*models.Module `json:"modules"`
	Count   int              `json:"count"`
}

type ListLogLevelsResponse struct {
	*Response
	LogLevels []*models.LogLevel `json:"log_levels"`
	Count     int                `json:"count"`
}
