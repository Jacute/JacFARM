package slogmultihandler

import (
	"context"
	"errors"
	"log/slog"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func New(handlers []slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Считаем включённым, если хотя бы один хендлер разрешает этот уровень.
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	errs := make([]error, 0, len(m.handlers))
	for _, h := range m.handlers {
		err := h.Handle(ctx, r)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{newHandlers}
}
