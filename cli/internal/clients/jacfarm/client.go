package jacfarm_client

import (
	"bytes"
	"cli_exploit_runner/pkg/common_config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrFlagFormatNotFound = errors.New("flag format not found")
)

const defaultPort = 15050
const defaultTimeout = 5 * time.Second

type Client struct {
	addr       string
	httpClient *http.Client
	token      string
}

type options struct {
	port    *int
	timeout time.Duration // int seconds
}

type Option func(opts *options) error

func New(host string, token string, opts ...Option) (*Client, error) {
	options := &options{}
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}

	// default params
	port := defaultPort
	timeout := defaultTimeout

	// apply options if use
	if options.port != nil {
		port = *options.port
	}
	if options.timeout != 0 {
		timeout = options.timeout
	}

	c := &Client{
		addr: fmt.Sprintf("%s:%d", host, port),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		token: token,
	}

	// trying to ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) error {
		opts.timeout = timeout
		return nil
	}
}

func WithCustomPort(port int) Option {
	return func(opts *options) error {
		p := port
		opts.port = &p
		return nil
	}
}

func (c *Client) GetTeams(ctx context.Context) ([]*Team, error) {
	url := fmt.Sprintf("http://%s/api/v1/service/teams", c.addr)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var teams []*Team
	if err := json.Unmarshal(data, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (c *Client) getConfig(ctx context.Context) ([]*Config, error) {
	url := fmt.Sprintf("http://%s/api/v1/service/config", c.addr)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cfg []*Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Client) GetFlagFormat(ctx context.Context) (string, error) {
	cfg, err := c.getConfig(ctx)
	if err != nil {
		return "", err
	}

	for _, cfgRow := range cfg {
		if cfgRow.Name == common_config.ConfigFlagFormatKey {
			return cfgRow.Value, nil
		}
	}

	return "", ErrFlagFormatNotFound
}

func (c *Client) SendFlags(ctx context.Context, flags []*ServiceFlag) error {
	url := fmt.Sprintf("http://%s/api/v1/service/flags", c.addr)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")

	data, err := json.Marshal(flags)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(data))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	return nil
}
