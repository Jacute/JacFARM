package main

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"fmt"
	"net"
	"strings"
)

const bufSize = 64

type Client struct {
	url   string
	token string
}

type FlagInfo struct {
	Flag string `json:"flag"`
	Msg  string `json:"msg"`
}

var NewClient plugins.NewClientFunc = func(url, token string) plugins.IClient {
	return &Client{
		url:   url,
		token: token,
	}
}

func (c *Client) SendFlags(ctx context.Context, flags []string) (map[string]*plugins.FlagResult, error) {
	conn, err := net.Dial("tcp", c.url)
	if err != nil {
		return nil, fmt.Errorf("error connecting: %w", err)
	}
	defer conn.Close()

	flagMap := make(map[string]*plugins.FlagResult)
	for _, flag := range flags {
		conn.Write([]byte(flag + "\n"))
		buf := make([]byte, bufSize)
		n, _ := conn.Read(buf)
		buf = buf[:n]

		status := models.FlagStatusReject
		if strings.Contains(string(buf), "Expired") {
			status = models.FlagStatusOld
		} else if strings.Contains(string(buf), "[OK]") {
			status = models.FlagStatusSuccess
		}

		flagMap[flag] = &plugins.FlagResult{
			Msg:    string(buf),
			Status: status,
		}
	}

	return flagMap, nil
}
