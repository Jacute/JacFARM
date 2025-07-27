package models

type FlagStatus string

type Team struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type Flag struct {
	Value       string `json:"value"`
	ExploitName string `json:"exploit_name"`
	GetFrom     *Team  `json:"victim_team"`
}
