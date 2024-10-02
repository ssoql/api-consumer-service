package interfaces

import (
	"context"

	"api-consumer-service/internal/dto"
)

type PostsDispatcher interface {
	Dispatch(ctx context.Context, postChan <-chan []dto.Post) error
}
