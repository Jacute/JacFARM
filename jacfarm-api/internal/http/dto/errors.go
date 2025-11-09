package dto

import "errors"

var (
	ErrLimitIncorrectType            = errors.New("limit should be number")
	ErrPageIncorrectType             = errors.New("page should be number")
	ErrLimitNegative                 = errors.New("limit should be positive number")
	ErrPageNegative                  = errors.New("page should be positive number")
	ErrLongExploitName               = errors.New("name should be shorter than 100 characters")
	ErrIncorrectFieldType            = errors.New("incorrect field 'type'")
	ErrRequirementsIncorrectMIMEType = errors.New("'requirements' has incorrect mime type")
	ErrExploitNotMatchType           = errors.New("'exploit' doesn't match 'type'")
	ErrTeamIdIncorrectType           = errors.New("team_id should be number")
	ErrStatusIdIncorrectType         = errors.New("status_id should be number")
	ErrIdIncorrectType               = errors.New("id should be number")
	ErrIdShouldBePos                 = errors.New("id should be positive")
)
