package jacfarm_client

import (
	"fmt"
	"net/http"
	"time"
)

const defaultPort = 15050
const defaultTimeout = 5 * time.Second

type Client struct {
	addr       string
	httpClient *http.Client
}

type options struct {
	port    *int
	timeout time.Duration // int seconds
}

type Option func(opts *options) error

func New(host string, opts ...Option) (*Client, error) {
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

	return &Client{
		addr: fmt.Sprintf("%s:%d", host, port),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
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
