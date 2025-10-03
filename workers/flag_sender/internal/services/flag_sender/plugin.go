package flag_sender

import (
	"flag_sender/pkg/plugins"
	"fmt"
	"plugin"
)

func loadPlugin(pluginPath, juryFlagURL, token string) (plugins.IClient, error) {
	sendPlugin, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("error opening send plugin: %w", err)
	}
	symbol, err := sendPlugin.Lookup("NewClient")
	if err != nil {
		return nil, fmt.Errorf("error looking up send plugin client: %w", err)
	}
	clientInit, ok := symbol.(*plugins.NewClientFunc)
	if !ok {
		return nil, fmt.Errorf("plugin client constructor not func(url, token string) Client")
	}
	pluginClient := (*clientInit)(juryFlagURL, token)

	return pluginClient, nil
}
