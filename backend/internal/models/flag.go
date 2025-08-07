package models

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
	CreatedAt         int64      `json:"created_at"` // ISO 8601 format
}

type FlagEnrich struct {
	ID                int64      `json:"id"`
	Value             string     `json:"value"`
	Status            FlagStatus `json:"status"`
	Exploit           *Exploit   `json:"exploit"`
	GetFrom           *Team      `json:"victim_team"`
	MessageFromServer string     `json:"message_from_server"`
	CreatedAt         string     `json:"created_at"` // ISO 8601 format
}
