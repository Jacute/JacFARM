package models

type FlagStatus string

type Flag struct {
	Value       string `json:"value"`
	ExploitName string `json:"exploit_name"`
	GetFrom     *Team  `json:"victim_team"`
}
