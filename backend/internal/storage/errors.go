package storage

import "errors"

var (
	ErrTeamAlreadyExists = errors.New("team with this ip already exists")
)
