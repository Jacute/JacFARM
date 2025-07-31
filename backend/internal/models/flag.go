package models

type FlagStatus string

const (
	FlagStatusPending    FlagStatus = "pending"
	FlagStatusInProgress FlagStatus = "in_progress"
	FlagStatusCompleted  FlagStatus = "completed"
	FlagStatusFailed     FlagStatus = "failed"
)

type Flag struct {
	ID                string     `json:"id"`
	Value             string     `json:"value"`
	Status            FlagStatus `json:"status"`
	Exploit           *Exploit   `json:"exploit"`
	GetFrom           *Team      `json:"victim_team"`
	MessageFromServer string     `json:"message_from_server"`
}
