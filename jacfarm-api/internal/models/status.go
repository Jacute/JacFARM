package models

type FlagStatus string

const (
	FlagStatusPending FlagStatus = "PENDING"
	FlagStatusOld     FlagStatus = "OLD"
	FlagStatusSuccess FlagStatus = "SUCCESS"
	FlagStatusReject  FlagStatus = "REJECT"
)

type Status struct {
	ID   int64      `json:"id"`
	Name FlagStatus `json:"name"`
}
