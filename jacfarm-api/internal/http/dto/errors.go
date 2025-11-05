package dto

import "fmt"

var (
	ErrLimitIncorrectType = fmt.Errorf("limit should be number")
	ErrPageIncorrectType  = fmt.Errorf("page should be number")
	ErrLimitNegative      = fmt.Errorf("limit should be positive number")
	ErrPageNegative       = fmt.Errorf("page should be positive number")
)
