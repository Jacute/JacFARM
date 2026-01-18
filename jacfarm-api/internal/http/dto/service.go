package dto

type ServicePutFlagRequest struct {
	Flags []*ServiceFlag `json:"flags"`
}

type ServiceFlag struct {
	Flag   string `json:"flag"`
	TeamID int64  `json:"team_id"`
}
