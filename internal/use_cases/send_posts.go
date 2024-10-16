package use_cases

import (
	"context"
	"fmt"
	"log"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type SendPostsUseCase struct {
	dispatcher interfaces.PostsDispatcher
}

func NewSendPostsUseCase(dispatcher interfaces.PostsDispatcher) *SendPostsUseCase {
	return &SendPostsUseCase{
		dispatcher: dispatcher,
	}
}

func (u *SendPostsUseCase) Send(ctx context.Context, ch <-chan []dto.Post) error {
	err := u.dispatcher.Dispatch(ctx, ch)
	if err != nil {
		log.Println("error: " + err.Error())

		return fmt.Errorf("failed to send posts: %w", err)
	}

	return nil
}
