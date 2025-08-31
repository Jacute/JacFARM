package storage

import "errors"

var (
	ErrTeamAlreadyExists   = errors.New("team with this ip already exists")
	ErrTeamNotFound        = errors.New("team not found")
	ErrConfigParamNotFound = errors.New("config parameter not found")
	ErrFlagNotUpdated      = errors.New("flag was not updated")
	ErrFlagAlreadyExists   = errors.New("flag with this value already exists")
	ErrExploitNotFound     = errors.New("exploit not found")
)
