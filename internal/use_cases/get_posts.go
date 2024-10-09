package use_cases

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type PostsGetter interface {
	Handle(ctx context.Context, url string) ([]dto.Post, error)
}

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

func (u *GetPostsUseCase) Handle(ctx context.Context, url string) ([]dto.Post, error) {
	var posts dto.Posts
	getData := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		callback := func() error {
			return u.client.DoRequest(req, &posts)
		}

		err = u.retryStrategy.Retry(callback)
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

func (u *GetPostsUseCase) Handlex(ctx context.Context, wg *sync.WaitGroup, jobsCounter int, pageChan <-chan string,
	postsChan chan<- []dto.Post) {
	for i := 0; i < jobsCounter; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for page := range pageChan {
				fmt.Printf("Worker %d fetching page %s\n", id, page)
				start := time.Now()

				posts, err := u.getPosts(ctx, page)
				if err != nil {
					fmt.Printf("Error fetching posts: %s\n", err)
					if errors.Is(err, context.Canceled) {
						return
					}
				}
				elapsed := time.Since(start)

				postsChan <- posts
				fmt.Printf("Worker %d finished page %s in %s\n", id, page, elapsed.String())
				time.Sleep(10 * time.Millisecond)
				//time.Sleep(3 * time.Second)
			}
		}(i)
	}
}

func (u *GetPostsUseCase) getPosts(ctx context.Context, url string) ([]dto.Post, error) {
	var posts dto.Posts
	getData := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		callback := func() error {
			return u.client.DoRequest(req, &posts)
		}

		err = u.retryStrategy.Retry(callback)
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
