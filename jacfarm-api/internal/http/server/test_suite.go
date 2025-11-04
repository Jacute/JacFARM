package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"log/slog"
	"net/url"
	"testing"
	"time"

	"github.com/jacute/prettylogger"
)

type TestSuite struct {
	app *HTTPServer
}

func NewTestSuite(t *testing.T, s handlers.Service) *TestSuite {
	h := handlers.New(s)
	log := slog.New(prettylogger.NewDiscardHandler())
	cfg := &config.HTTPConfig{
		Host:         "localhost",
		Port:         8080,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		CORS: &config.CORSConfig{
			AllowedOrigins: []string{"http://test"},
		},
	}
	app := New(
		log,
		cfg,
		"test",
		h,
	)
	go app.Start()
	t.Cleanup(func() {
		app.Stop()
	})

	return &TestSuite{
		app: app,
	}
}

func queryParamsToString(queryParams map[string][]string) string {
	q := url.Values{}
	for k, vs := range queryParams {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	return q.Encode()
}
