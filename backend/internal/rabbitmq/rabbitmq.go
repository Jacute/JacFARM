package rabbitmq

import (
	"JacFARM/internal/config"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const flagsQueueName = "flags"

type Rabbit struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	flagsQueue *amqp.Queue
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
	ch, err := conn.Channel()
	if err != nil {
		panic("failed to create rabbitmq channel: " + err.Error())
	}
	q, err := ch.QueueDeclare(
		flagsQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		panic("failed to declare queue: " + err.Error())
	}
	return &Rabbit{
		conn:       conn,
		ch:         ch,
		flagsQueue: &q,
	}
}

func (r *Rabbit) Close() error {
	if err := r.ch.Close(); err != nil {
		return fmt.Errorf("error closing channel %v", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("error closing conn %v", err)
	}
	return nil
}
