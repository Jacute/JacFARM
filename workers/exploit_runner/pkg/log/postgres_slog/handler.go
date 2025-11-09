package postgresslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/jacute/prettylogger/handlers"
)

type LogStorage interface {
	WriteLog(
		ctx context.Context,
		module, op, level, value, exploitId string,
		attrs map[string]any,
		createdAt time.Time,
	) error
}

type PostgresHandler struct {
	module string
	db     LogStorage
	H      slog.Handler
	B      *bytes.Buffer
	M      *sync.Mutex
	Opts   *slog.HandlerOptions
	W      io.Writer
}

func (h *PostgresHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Opts.Level.Level()
}

func (h *PostgresHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PostgresHandler{
		H:      h.H.WithAttrs(attrs),
		B:      h.B,
		M:      h.M,
		Opts:   h.Opts,
		W:      h.W,
		module: h.module,
		db:     h.db,
	}
}

func (h *PostgresHandler) WithGroup(name string) slog.Handler {
	return &PostgresHandler{
		H:      h.H.WithGroup(name),
		B:      h.B,
		M:      h.M,
		Opts:   h.Opts,
		W:      h.W,
		module: h.module,
		db:     h.db,
	}
}

func (h *PostgresHandler) computeAttr(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.M.Lock()
	defer func() {
		h.B.Reset()
		h.M.Unlock()
	}()
	if err := h.H.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %W", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.B.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshalling inner handler's Handle attrs: %W", err)
	}
	return attrs, nil
}

func (h *PostgresHandler) Handle(ctx context.Context, r slog.Record) error {
	// get attrs
	level := r.Level.String()
	createdAt := r.Time.UTC()

	attrs, err := h.computeAttr(ctx, r)
	if err != nil {
		return err
	}

	// add db fields
	var operation string
	if op, ok := attrs["op"]; ok {
		opStr, ok := op.(string)
		if ok {
			operation = opStr
			delete(attrs, "op")
		}
	} else {
		return nil // write in database only rows with attr "op"
	}

	var exploitId string
	if id, ok := attrs["exploit_id"]; ok {
		idStr, ok := id.(string)
		if ok {
			exploitId = idStr
			delete(attrs, "exploit_id")
		}
	}

	err = h.db.WriteLog(
		ctx,
		h.module,
		operation,
		level,
		r.Message,
		exploitId,
		attrs,
		createdAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewHandler(module string, db LogStorage, opts *slog.HandlerOptions) *PostgresHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	b := &bytes.Buffer{}
	return &PostgresHandler{
		module: module,
		db:     db,
		H: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: handlers.SupressDefaults(opts.ReplaceAttr),
		}),
		B:    b,
		M:    &sync.Mutex{},
		Opts: opts,
		W:    nil,
	}
}
