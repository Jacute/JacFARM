package rabbitmq

import (
	"JacFARM/pkg/rabbitmq_dto"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *Rabbit) PublishFlag(flag *rabbitmq_dto.Flag) error {
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
