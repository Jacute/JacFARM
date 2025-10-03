package rabbitmq

import (
	"JacFARM/pkg/rabbitmq_dto"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         output,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rabbit) GetFlagsCount() (int, error) {
	values := url.Values{
		"lengths_age":     []string{"60"},
		"lengths_incr":    []string{"5"},
		"msg_rates_age":   []string{"60"},
		"msg_rates_incr":  []string{"5"},
		"data_rates_age":  []string{"60"},
		"data_rates_incr": []string{"5"},
	}
	getFlagsURL := fmt.Sprintf(
		"http://%s:%d/api/queues/%%2f/%s?%s",
		r.cfg.ManagementHost,
		r.cfg.ManagementPort,
		flagsQueueName,
		values.Encode(),
	)

	req, err := http.NewRequest("GET", getFlagsURL, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(r.cfg.Username, r.cfg.Password)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	var data QueueInfo
	if err = sonic.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	return data.MessagesCount, nil
}
