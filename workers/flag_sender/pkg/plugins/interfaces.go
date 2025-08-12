package plugins

type IClient interface {
	SendFlags([]string) (map[string]*FlagResult, error)
}

type NewClientFunc func(url, token string) IClient
