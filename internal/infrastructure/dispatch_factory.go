package infrastructure

import (
	"errors"

	"api-consumer-service/internal/infrastructure/cli"
	"api-consumer-service/internal/infrastructure/queue"
	"api-consumer-service/internal/use_cases/interfaces"
)

func CreateDispatcher(envName string, queueUrl, queueName string) (interfaces.PostsDispatcher, error) {
	switch envName {
	case "local":
	case "dev":
		return cli.NewCliDispatch(), nil
	case "prod":
		client, err := queue.NewRabbitMqClient(queueUrl, queueName)
		if err != nil {
			return nil, err
		}
		return queue.NewRabbitmqDispatcher(client), nil
	}

	return nil, errors.New("invalid dispatch type")
}
