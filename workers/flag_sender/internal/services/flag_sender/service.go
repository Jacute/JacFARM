package flag_sender

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/jacute/prettylogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockgen -source=service.go -destination=./mocks/storage_mock.go -package=mocks -mock_names=storage=StorageMock .
type storage interface {
	PutFlags(ctx context.Context, flags []*models.Flag) ([]int64, error)
	GetConfigParameter(ctx context.Context, key string) (string, error)
	GetFlagValuesByStatus(ctx context.Context, status models.FlagStatus) ([]string, error)
	UpdateStatusForOldFlags(ctx context.Context, flagTTL time.Duration) (int64, error)
	UpdateFlagByResult(ctx context.Context, flag string, result *plugins.FlagResult) error
}

//go:generate mockgen -source=service.go -destination=./mocks/queue_mock.go -package=mocks -mock_names=queue=QueueMock .
type queue interface {
	GetFlagChan() (<-chan amqp.Delivery, error)
}

type FlagSender struct {
	log          *slog.Logger
	queue        queue
	db           storage
	cfg          *config
	pluginClient plugins.IClient
	stopChan     chan struct{}
}

func New(log *slog.Logger, pluginDir string, q queue, db storage) (*FlagSender, error) {
	fs := &FlagSender{
		log:      log,
		queue:    q,
		db:       db,
		stopChan: make(chan struct{}),
		cfg:      &config{},
	}
	err := fs.loadConfig(context.Background(), true)
	if err != nil {
		return nil, fmt.Errorf("error loading flag sender config: %w", err)
	}

	// load plugin
	pluginPath := path.Join(pluginDir, fs.cfg.plugin+".so")
	plugin, err := loadPlugin(pluginPath, fs.cfg.juryFlagURL, fs.cfg.token)
	if err != nil {
		return nil, fmt.Errorf("error loading plugin: %w", err)
	}
	fs.pluginClient = plugin

	return fs, nil
}

func (fs *FlagSender) Start() error {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("Starting FlagSender service")

	flagChan, err := fs.queue.GetFlagChan()
	if err != nil {
		log.Error("Failed to get flag channel", prettylogger.Err(err))
		return err
	}

	batch := make([]amqp.Delivery, 0, fs.cfg.submitLimit)
	for {
		timer := time.NewTimer(fs.cfg.submitPeriod)
		select {
		case flag, ok := <-flagChan:
			if !ok {
				log.Info("flag channel closed")
				return nil
			}

			batch = append(batch, flag)
			if len(batch) >= fs.cfg.submitLimit {
				sendCtx, cancel := context.WithTimeout(context.Background(), fs.cfg.submitTimeout)
				last := batch[len(batch)-1]
				if err := fs.processBatch(sendCtx, batch); err != nil {
					log.Error("failed to process flag", prettylogger.Err(err))
					last.Nack(true, true) // if error, requeue the message
					batch = batch[:0]
					cancel()
					continue
				}

				last.Ack(true)
				batch = batch[:0]
				cancel()
			}
		case <-timer.C:
			sendCtx, cancel := context.WithTimeout(context.Background(), fs.cfg.submitTimeout)
			err := fs.loadConfig(sendCtx, false)
			if err != nil {
				log.Error("error reloading flag sender config from db", prettylogger.Err(err))
				cancel()
				continue
			}
			if len(batch) > 0 {
				last := batch[len(batch)-1]
				if err := fs.processBatch(sendCtx, batch); err != nil {
					log.Error("failed to process flag", prettylogger.Err(err))
					last.Nack(true, true) // if error, requeue the message
					cancel()
					continue
				}

				last.Ack(true)
				batch = batch[:0]
			}
			cancel()
		case <-fs.stopChan:
			log.Info("flag sender stopped")
			timer.Stop()
			return nil
		}
	}
}

func (fs *FlagSender) Stop() {
	const op = "service.flag_sender.Stop"
	log := fs.log.With(slog.String("op", op))

	log.Debug("stopping flag saver service")
	close(fs.stopChan)
}
