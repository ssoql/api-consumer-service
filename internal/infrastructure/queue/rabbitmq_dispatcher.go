package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type RabbitmqDispatcher struct {
	client *RabbitMqClient
}

func NewRabbitmqDispatcher(client *RabbitMqClient) interfaces.PostsDispatcher {
	return &RabbitmqDispatcher{
		client: client,
	}
}

func (r *RabbitmqDispatcher) Dispatch(ctx context.Context, postChan <-chan []dto.Post) error {
	defer r.client.Close()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case posts, ok := <-postChan:
			if !ok {
				return nil
			}

			for _, post := range posts {
				err := r.publishPost(post)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (r *RabbitmqDispatcher) publishPost(post dto.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %w", err)
	}
	err = r.client.channel.Publish(
		"",                  // exchange
		r.client.queue.Name, // routing key (queue name)
		false,               // mandatory
		false,               // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish post: %w", err)
	}

	return nil
}
