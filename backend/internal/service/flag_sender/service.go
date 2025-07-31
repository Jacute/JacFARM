package flag_sender

import (
	"JacFARM/internal/models"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type storage interface {
	PutFlag(flag *models.Flag) error
}

type queue interface {
	GetFlagChan() (<-chan amqp.Delivery, error)
}

type FlagSender struct {
	log      *slog.Logger
	queue    queue
	db       storage
	stopChan chan struct{}
}

func New(log *slog.Logger, q queue, db storage) *FlagSender {
	return &FlagSender{
		log:      log,
		queue:    q,
		db:       db,
		stopChan: make(chan struct{}),
	}
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

	go func() {
		for {
			select {
			case flag, ok := <-flagChan:
				if !ok {
					log.Info("Flag channel closed")
					return
				}

				fs.log.Info("Received flag", "body", string(flag.Body))

				if err := fs.processFlag(flag.Body); err != nil {
					log.Error("Failed to process flag", "error", err)
					flag.Nack(false, true) // if error, requeue the message
					continue
				}

				flag.Ack(false) // if success, send ack
			case <-fs.stopChan:
				return
			}
		}
	}()

	return nil
}

func (fs *FlagSender) processFlag(flag []byte) error {
	const op = "service.flag_sender.processFlag"
	log := fs.log.With(slog.String("op", op))
	log.Info("Processing flag", slog.String("flag", string(flag)))

	time.Sleep(500 * time.Millisecond) // process delay simulation
	return nil
}

func (fs *FlagSender) Stop() {
	const op = "service.flag_sender.Stop"

	close(fs.stopChan)
	fs.log.Info("FlagSender stopped gracefully", slog.String("op", op))
}
