package flag_sender

import (
	"context"
	"log/slog"
)

type storage interface {
}

type FlagSender struct {
	log *slog.Logger
	db  storage
	cfg *config
}

type config struct {
}

func New(log *slog.Logger, db storage) *FlagSender {
	er := &FlagSender{
		log: log,
		db:  db,
	}
	err := er.loadConfig(context.Background())
	if err != nil {
		panic("error loading exploit runner config: " + err.Error())
	}

	return er
}

func (er *FlagSender) loadConfig(ctx context.Context) error {
	const op = "service.flag_sender.LoadConfig"
	log := er.log.With(slog.String("op", op))

	_ = log
	// TODO: imlement loadConfig

	return nil
}
