package cli

import (
	"context"
	"fmt"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type CliDispatch struct{}

func NewCliDispatch() interfaces.PostsDispatcher {
	return new(CliDispatch)
}

func (c *CliDispatch) Dispatch(ctx context.Context, postChan <-chan []dto.Post) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case posts, ok := <-postChan:
			if !ok {
				return nil
			}

			for _, post := range posts {
				fmt.Printf("Post %d: %s\n", post.ID, post.Title)
			}
		}
	}

	return nil
}
