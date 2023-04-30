package rabbbitmq

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type rabbbitmq struct {
	conn *amqp.Connection
}

func (r *rabbbitmq) Subscribe(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error) {
	panic("not implemented") // TODO: Implement
}

func (r *rabbbitmq) Publish(topic, subTopic string, msg io.Reader) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "could not open channel connection")
	}
	if err := ch.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrapf(err, "could not declare exchange: %v", topic)
	}
	body, err := ioutil.ReadAll(msg)
	if err != nil {
		return errors.Wrap(err, "could not read message reader")
	}
	err = ch.Publish(
		topic,
		subTopic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "could not publish message: %v; %v", topic, subTopic)
	}
	return nil
}

func (r *rabbbitmq) Close() error {
	return r.conn.Close()
}
