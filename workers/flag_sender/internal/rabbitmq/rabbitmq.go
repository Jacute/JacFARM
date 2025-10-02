package rabbitmq

import (
	"flag_sender/internal/config"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const flagsQueueName = "flags"

type Rabbit struct {
	conn   *amqp.Connection
	readCh *amqp.Channel
}

func New(cfg *config.RabbitMQConfig) *Rabbit {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	))
	if err != nil {
		panic("failed to connect to RabbitMQ: " + err.Error())
	}
	readCh, err := conn.Channel()
	if err != nil {
		panic("failed to create rabbitmq channel: " + err.Error())
	}
	_, err = readCh.QueueDeclare(
		flagsQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			"x-message-deduplication": true,
		},
	)
	if err != nil {
		panic("failed to declare queue: " + err.Error())
	}
	return &Rabbit{
		conn:   conn,
		readCh: readCh,
	}
}

func (r *Rabbit) Close() error {
	if err := r.readCh.Close(); err != nil {
		return fmt.Errorf("error closing channel %v", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("error closing conn %v", err)
	}
	return nil
}
