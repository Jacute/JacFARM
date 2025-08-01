package flag_saver

import (
	"JacFARM/internal/models"
	"log/slog"
	"sync"

	"github.com/jacute/prettylogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type storage interface {
	PutFlag(flag *models.Flag) error
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
					log.Error("Failed to process flag", "error", err)
					flag.Nack(false, true) // if error, requeue the message
					return
				}

				flag.Ack(false) // if success, send ack
			}()
		case <-fs.stopChan:
			processFlagWg.Wait()
			return nil
		}
	}
}

func (fs *FlagSaver) Stop() {
	const op = "service.flag_sender.Stop"
	log := fs.log.With(slog.String("op", op))

	log.Info("stopping flag saver service")
	close(fs.stopChan)
}
