package flag_saver

import (
	"JacFARM/internal/models"
	"context"
	"errors"
	"log/slog"
	"sync"

	storage_errors "JacFARM/internal/storage"

	"github.com/jacute/prettylogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type storage interface {
	PutFlag(ctx context.Context, flag *models.Flag) (int64, error)
}

type queue interface {
	GetFlagChan() (<-chan amqp.Delivery, error)
}

type FlagSaver struct {
	log      *slog.Logger
	queue    queue
	db       storage
	stopChan chan struct{}
}

func New(log *slog.Logger, q queue, db storage) *FlagSaver {
	return &FlagSaver{
		log:      log,
		queue:    q,
		db:       db,
		stopChan: make(chan struct{}),
	}
}

func (fs *FlagSaver) Start() error {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("Starting FlagSender service")

	flagChan, err := fs.queue.GetFlagChan()
	if err != nil {
		log.Error("Failed to get flag channel", prettylogger.Err(err))
		return err
	}

	var processFlagWg sync.WaitGroup
	for {
		select {
		case flag, ok := <-flagChan:
			if !ok {
				processFlagWg.Wait()
				log.Info("Flag channel closed")
				return nil
			}

			log.Info("Received flag", "body", string(flag.Body))

			processFlagWg.Add(1)
			go func() {
				defer processFlagWg.Done()
				if err := fs.processFlag(flag.Body); err != nil {
					if errors.Is(err, storage_errors.ErrFlagAlreadyExists) {
						log.Warn("Flag already exists, skipping", slog.String("flag", string(flag.Body)))
						flag.Ack(false) // if flag already exists, send ack
						return
					}
					log.Error("Failed to process flag", prettylogger.Err(err))
					flag.Nack(false, true) // if error, requeue the message
					return
				}

				flag.Ack(false) // if success, send ack
			}()
		case <-fs.stopChan:
			processFlagWg.Wait()
			log.Info("flag saver stopped")
			return nil
		}
	}
}

func (fs *FlagSaver) Stop() {
	const op = "service.flag_sender.Stop"
	log := fs.log.With(slog.String("op", op))

	log.Debug("stopping flag saver service")
	close(fs.stopChan)
}
