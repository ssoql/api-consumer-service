package queue

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMqClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewRabbitMqClient(amqpURL, queueName string) (*RabbitMqClient, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close rabbitmq connection: %v\n", err)
		}
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
		close(conn, ch)
		return nil, err
	}

	return &RabbitMqClient{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMqClient) Close() {
	close(r.conn, r.channel)
}

func close(connection *amqp091.Connection, channel *amqp091.Channel) {
	if err := connection.Close(); err != nil {
		log.Printf("failed to close rabbitmq connection: %v\n", err)
	}
	if err := channel.Close(); err != nil {
		log.Printf("failed to close rabbitmq channel: %v\n", err)
	}
}
