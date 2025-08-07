package plugins

import "JacFARM/internal/models"

type FlagResult struct {
	Status models.FlagStatus `json:"status"`
	Msg    string            `json:"msg"`
}
