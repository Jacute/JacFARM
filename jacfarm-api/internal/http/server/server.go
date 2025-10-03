package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"context"
	"fmt"
	"log/slog"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/jacute/prettylogger"
)

const shutdownTimeout = 5 * time.Second

type HTTPServer struct {
	log    *slog.Logger
	cfg    *config.HTTPConfig
	router *fiber.App
}

func New(log *slog.Logger, cfg *config.HTTPConfig, apiKey string, h *handlers.Handlers) *HTTPServer {
	return &HTTPServer{
		log:    log,
		cfg:    cfg,
		router: setupRouter(h, cfg, apiKey),
	}
}

func (s *HTTPServer) Start() {
	const op = "server.Start"
	log := s.log.With(slog.String("op", op))

	log.Info("starting http server", slog.Int("port", s.cfg.Port))
	if err := s.router.Listen(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)); err != nil {
		log.Error("error starting http server", prettylogger.Err(err))
	}
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.router.ShutdownWithContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
