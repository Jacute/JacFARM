package rabbitmq_dto

import "time"

type FlagSourceType string

var (
	LocalExploitSourceType  FlagSourceType = "LOCAL_EXPLOIT"
	FarmExploitSourceType   FlagSourceType = "FARM_EXPLOIT"
	ManualSendingSourceType FlagSourceType = "MANUAL_SENDING"
)

type Flag struct {
	Value      string         `json:"value"`
	ExploitID  string         `json:"exploit"`
	TeamID     int64          `json:"victim_team"`
	SourceType FlagSourceType `json:"source_type"`
	CreatedAt  time.Time      `json:"created_at"`
}

type QueueResponse struct {
}
