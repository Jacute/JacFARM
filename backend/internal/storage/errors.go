package storage

import "errors"

var (
	ErrTeamAlreadyExists        = errors.New("team with this ip already exists")
	ErrConfigParamAlreadyExists = errors.New("config parameter with this key already exists")
)
