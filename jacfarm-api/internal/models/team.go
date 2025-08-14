package models

import "net"

type Team struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	IP   net.IP `json:"ip"`
}
