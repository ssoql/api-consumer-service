package queue

import (
	"github.com/streadway/amqp"
)

type RabbitMqClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMqClient(amqpURL, queueName string) (*RabbitMqClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // todo handle error
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()   // todo handle error
		conn.Close() // todo handle error
		return nil, err
	}

	return &RabbitMqClient{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMqClient) Close() {
	r.conn.Close()    // todo handle error
	r.channel.Close() // todo handle error
}
