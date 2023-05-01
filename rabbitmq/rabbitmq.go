package rabbbitmq

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type rabbbitmq struct {
	conn *amqp.Connection
}

// New creates new Message Broker connection
func New(user, password, host, port string) (*rabbbitmq, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", user, password, host, port))
	if err != nil {
		return nil, err
	}
	return &rabbbitmq{conn}, nil
}

func (r *rabbbitmq) Subscribe(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error) error {
	// Open new channel
	ch, err := r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "could not create a channel connection")
	}
	defer ch.Close()

	go func() {
		<-stop
		ch.Close()
	}()
	queueName := clientID
	// Queue Dleclaration
	if _, err := ch.QueueDeclare(
		queueName, // QueueName
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrapf(err, "could not declare a queue: %v", clientID)
	}
	// Exchange and Queue Bind Declarations
	for exchange, routeKeys := range topics {
		if err := ch.ExchangeDeclare(
			exchange,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return errors.Wrapf(err, "could not declare exchange: %v", exchange)
		}
		for _, key := range routeKeys {
			if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
				return errors.Wrapf(err, "could not bind queue: %v; route key: %v; exchange: %v", queueName, key, exchange)
			}
		}
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "could not get consume channel: %v", queueName)
	}
	for d := range msgs {
		if err := callback(bytes.NewBuffer(d.Body)); err != nil {
			d.Nack(false, true)
			return err
		}
		d.Ack(false)
	}
	return nil
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
