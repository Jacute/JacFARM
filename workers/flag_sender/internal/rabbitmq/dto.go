package rabbitmq

type Flag struct {
	Value     string `json:"value"`
	ExploitID string `json:"exploit"`
	TeamID    int64  `json:"victim_team"`
	IsLocal   bool   `json:"is_local"`
}
