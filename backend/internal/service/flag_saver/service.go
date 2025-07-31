package flag_saver

import (
	"JacFARM/internal/models"
	"log/slog"

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
