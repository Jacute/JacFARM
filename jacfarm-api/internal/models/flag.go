package models

import (
	"time"
)

type FlagStatus string

const (
	FlagStatusPending FlagStatus = "PENDING"
	FlagStatusOld     FlagStatus = "OLD"
	FlagStatusSuccess FlagStatus = "SUCCESS"
	FlagStatusReject  FlagStatus = "REJECT"
)

type Flag struct {
	ID                int64      `json:"id"`
	Value             string     `json:"value"`
	Status            FlagStatus `json:"status"`
	ExploitID         string     `json:"exploit_id"`
	GetFrom           int64      `json:"victim_team_id"`
	MessageFromServer string     `json:"message_from_server"`
	CreatedAt         time.Time  `json:"created_at"`
}

type FlagEnrich struct {
	ID                int64      `json:"id"`
	Value             string     `json:"value"`
	Status            FlagStatus `json:"status"`
	ExploitName       string     `json:"exploit_name"`
	VictimIP          string     `json:"victim_team"`
	MessageFromServer string     `json:"message_from_server"`
	CreatedAt         time.Time  `json:"created_at"`
}
