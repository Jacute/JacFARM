package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *Rabbit) GetFlagChan() (<-chan amqp.Delivery, error) {
	msgs, err := r.readCh.Consume(
		r.flagsQueue.Name,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
