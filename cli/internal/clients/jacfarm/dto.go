package jacfarm_client

import "net"

type Config struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

type Team struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IP   net.IP `json:"ip"`
}

type ListTeamsResponse struct {
	*Response
	Teams []*Team `json:"teams"`
	Count int     `json:"count"`
}

type GetConfigResponse struct {
	*Response
	Config []*Config `json:"config"`
	Count  int       `json:"count"`
}

type ServicePutFlagRequest struct {
	Flags []*ServiceFlag `json:"flags"`
}

type ServiceFlag struct {
	Flag   string `json:"flag"`
	TeamID int64  `json:"team_id"`
}
