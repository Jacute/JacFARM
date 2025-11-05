package dto

import "fmt"

var (
	ErrLimitIncorrectType            = fmt.Errorf("limit should be number")
	ErrPageIncorrectType             = fmt.Errorf("page should be number")
	ErrLimitNegative                 = fmt.Errorf("limit should be positive number")
	ErrPageNegative                  = fmt.Errorf("page should be positive number")
	ErrLongExploitName               = fmt.Errorf("name should be shorter than 100 characters")
	ErrIncorrectFieldType            = fmt.Errorf("incorrect field 'type'")
	ErrRequirementsIncorrectMIMEType = fmt.Errorf("'requirements' has incorrect mime type")
	ErrExploitNotMatchType           = fmt.Errorf("'exploit' doesn't match 'type'")
)
