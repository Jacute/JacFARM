package flag_sender

import (
	"context"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger"
)

type storage interface {
	GetConfigParameter(ctx context.Context, key string) (string, error)
}

type FlagSender struct {
	log          *slog.Logger
	db           storage
	cfg          *config
	shutdownChan chan struct{}
}

func New(log *slog.Logger, db storage) *FlagSender {
	er := &FlagSender{
		log:          log,
		db:           db,
		cfg:          &config{},
		shutdownChan: make(chan struct{}),
	}
	err := er.loadConfig(context.Background())
	if err != nil {
		panic("error loading exploit runner config: " + err.Error())
	}

	return er
}

func (fs *FlagSender) Start() {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("starting flag sender service")

	ticker := time.NewTicker(fs.cfg.submitPeriod)
	for {
		select {
		case <-ticker.C:
			// every tick load new config from db
			err := fs.loadConfig(context.Background())
			if err != nil {
				log.Error(
					"error reloading flag sender config from db",
					prettylogger.Err(err),
				)
				continue
			}
			// TODO: implement flag sender logic
		case <-fs.shutdownChan:
			log.Info("shutting down flag sender service")
			return
		}
	}
}

func (fs *FlagSender) Stop() {
	const op = "service.flag_sender.Stop"
	log := fs.log.With(slog.String("op", op))

	log.Info("stopping flag sender service")
	fs.shutdownChan <- struct{}{}
}
