package dto

import "errors"

var (
	ErrLimitIncorrect = errors.New("limit should be number")
	ErrPageIncorrect  = errors.New("page should be number")

	ErrLimitNegative = errors.New("limit should be positive number")
	ErrPageNegative  = errors.New("page should be positive number")
)
