package flag_sender

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"fmt"
	"log/slog"
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
	pluginDir    string
	pluginClient plugins.IClient
	stopChan     chan struct{}
}

func New(log *slog.Logger, pluginDir string, q queue, db storage) (*FlagSender, error) {
	fs := &FlagSender{
		log:       log,
		queue:     q,
		db:        db,
		pluginDir: pluginDir,
		stopChan:  make(chan struct{}),
		cfg:       &config{},
	}
	err := fs.loadConfig(context.Background(), true)
	if err != nil {
		return nil, fmt.Errorf("error loading flag sender config: %w", err)
	}

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
				continue
			}

			if len(batch) < fs.cfg.submitLimit {
				batch = append(batch, flag)
			} else {
				flag.Nack(false, true) // batch is full, requeue
			}
		case <-timer.C:
			sendCtx, cancel := context.WithTimeout(context.Background(), fs.cfg.submitTimeout)
			err := fs.loadConfig(sendCtx, false)
			if err != nil {
				log.Error("error reloading flag sender config from db", prettylogger.Err(err))
				continue
			}

			if len(batch) > 0 {
				err := fs.processBatch(sendCtx, batch)
				if err != nil {
					log.Error("failed to process flag", prettylogger.Err(err))
				}
				for _, msg := range batch {
					if err != nil {
						// requeue every msg
						msg.Nack(false, true)
					} else {
						msg.Ack(false)
					}
				}
				batch = batch[:0]
			} else {
				log.Info("no flags to process")
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

	log.Debug("stopping flag sender service")
	close(fs.stopChan)
}
