package main

import (
	"JacFARM/internal/models"
	"JacFARM/pkg/plugin_interfaces"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

type Client struct {
	url   string
	token string
}

type FlagInfo struct {
	Flag string `json:"flag"`
	Msg  string `json:"msg"`
}

var NewClient plugin_interfaces.NewClientFunc = func(url, token string) plugin_interfaces.IClient {
	return &Client{
		url:   url,
		token: token,
	}
}

func (c *Client) SendFlags(flags []string) (map[string]models.FlagStatus, error) {
	data, err := sonic.Marshal(flags)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Team-Token", c.token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("incorrect status code: %d, body: %s", res.StatusCode, body)
	}

	var flagInfos []*FlagInfo
	if err := sonic.Unmarshal(body, &flagInfos); err != nil {
		return nil, err
	}
	flagMap := make(map[string]models.FlagStatus)
	for _, flagInfo := range flagInfos {
		flagStatus := models.FlagStatusReject
		if strings.Contains(flagInfo.Msg, "accepted") {
			flagStatus = models.FlagStatusSuccess
		}
		if strings.Contains(flagInfo.Msg, "old") {
			flagStatus = models.FlagStatusOld
		}
		flagMap[flagInfo.Flag] = flagStatus
	}

	return flagMap, nil
}
