package server

import (
	"JacFARM/internal/config"
	"JacFARM/internal/http/handlers"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/jacute/prettylogger"
)

const shutdownTimeout = 5 * time.Second

type HTTPServer struct {
	log         *slog.Logger
	cfg         *config.HTTPConfig
	router      *fiber.App
	pprofRouter *http.Server
}

func New(log *slog.Logger, cfg *config.HTTPConfig, apiKey string, h *handlers.Handlers) *HTTPServer {
	pprofServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.PprofPort),
		Handler: setupPprofRouter(),
	}
	return &HTTPServer{
		log:         log,
		cfg:         cfg,
		router:      setupRouter(h, cfg, apiKey),
		pprofRouter: pprofServer,
	}
}

func (s *HTTPServer) Test(req *http.Request) (*http.Response, error) {
	return s.router.Test(req)
}

func (s *HTTPServer) Start() {
	const op = "server.Start"
	log := s.log.With(slog.String("op", op))

	log.Info("starting http server", slog.Int("port", s.cfg.Port))
	go func() {
		if err := s.pprofRouter.ListenAndServe(); err != nil {
			log.Error("error starting http pprof server", prettylogger.Err(err))
		}
	}()
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

	err = s.pprofRouter.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
