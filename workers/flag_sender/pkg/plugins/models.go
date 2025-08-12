package plugins

import "flag_sender/internal/models"

type FlagResult struct {
	Status models.FlagStatus `json:"status"`
	Msg    string            `json:"msg"`
}
