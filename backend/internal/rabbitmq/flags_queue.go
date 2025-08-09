package rabbitmq

import (
	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *Rabbit) PublishFlag(flag *Flag) error {
	output, err := sonic.Marshal(flag)
	if err != nil {
		return err
	}

	err = r.writeCh.Publish(
		"",
		r.flagsQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        output,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

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
