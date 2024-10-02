package use_cases

import (
	"context"
	"net/http"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type GetPostsUseCase struct {
	client        interfaces.ApiConsumer
	retryStrategy interfaces.Retryer
}

func NewGetPostsUseCase(client interfaces.ApiConsumer, retryStrategy interfaces.Retryer) *GetPostsUseCase {
	return &GetPostsUseCase{
		client:        client,
		retryStrategy: retryStrategy,
	}
}

func (u *GetPostsUseCase) GetPosts(ctx context.Context, url string) ([]dto.Post, error) {
	var posts dto.Posts
	getData := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		callback := func() error {
			return u.client.DoRequest(req, &posts)
		}

		err = u.retryStrategy.Retry(ctx, callback)
		if err != nil {
			return err
		}

		return nil
	}

	if err := getData(); err != nil {
		return nil, err
	}

	return posts.Posts, nil
}
