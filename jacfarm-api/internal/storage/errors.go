package storage

import "errors"

var (
	ErrTeamAlreadyExists        = errors.New("team with this ip already exists")
	ErrTeamNotFound             = errors.New("team not found")
	ErrConfigParamAlreadyExists = errors.New("config parameter with this key already exists")
	ErrFlagNotUpdated           = errors.New("flag was not updated")
	ErrFlagAlreadyExists        = errors.New("flag with this value already exists")
	ErrExploitNotFound          = errors.New("exploit not found")
)
