package plugins

import "context"

type IClient interface {
	SendFlags(ctx context.Context, flags []string) (map[string]*FlagResult, error)
}

type NewClientFunc func(url, token string) IClient
