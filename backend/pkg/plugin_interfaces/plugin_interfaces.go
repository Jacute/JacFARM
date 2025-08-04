package plugin_interfaces

import "JacFARM/internal/models"

type IClient interface {
	SendFlags([]string) (map[string]models.FlagStatus, error)
}

type NewClientFunc func(url, token string) IClient
