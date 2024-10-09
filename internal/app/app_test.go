package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases"
	"api-consumer-service/internal/use_cases/mocks"
)

type appTestSuite struct{}

func (a *appTestSuite) createPageChan(length int) chan string {
	ch := make(chan string, length)
	defer close(ch)

	for i := 1; i <= length; i++ {
		ch <- fmt.Sprintf("url_to_page_%d", i)
	}
	return ch
}

func (a *appTestSuite) useCaseMock() use_cases.PostsGetter {
	m := &mocks.MockPostsGetter{}
	m.On("Handle", mock.Anything, mock.Anything).Return([]dto.Post{
		{ID: 1, UserID: 1, Title: "TEst", Body: "nothing special"},
	}, nil)

	return m
}

func (a *appTestSuite) useCaseFailureMock() use_cases.PostsGetter {
	m := &mocks.MockPostsGetter{}
	m.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("API error"))

	return m
}

func (a *appTestSuite) useCaseContextCancelledMock() use_cases.PostsGetter {
	m := &mocks.MockPostsGetter{}
	m.On("Handle", mock.Anything, mock.Anything).Return(nil, context.Canceled)

	return m
}

func Test_fetchAllPosts(t *testing.T) {
	currentTest := &appTestSuite{}
	tests := []struct {
		name        string
		ctx         context.Context
		useCase     use_cases.PostsGetter
		jobsCounter int
		assertFunc  func(t *testing.T, pages chan string, posts chan []dto.Post)
	}{
		{
			name:        "success",
			ctx:         context.Background(),
			useCase:     currentTest.useCaseMock(),
			jobsCounter: 2,
			assertFunc: func(t *testing.T, pages chan string, posts chan []dto.Post) {
				assert.Equal(t, 0, len(pages), "pageChann must be empty")
				assert.NotEqual(t, 0, len(posts), "postsChann can't be empty")
			},
		},
		{
			name:        "failure",
			ctx:         context.Background(),
			useCase:     currentTest.useCaseFailureMock(),
			jobsCounter: 2,
			assertFunc: func(t *testing.T, pages chan string, posts chan []dto.Post) {
				assert.Equal(t, 0, len(pages), "pageChann must be empty")
				assert.Equal(t, 0, len(posts), "postsChann must be empty")
			},
		},
		{
			name:        "context-cancelled",
			ctx:         context.Background(),
			useCase:     currentTest.useCaseContextCancelledMock(),
			jobsCounter: 2,
			assertFunc: func(t *testing.T, pages chan string, posts chan []dto.Post) {
				assert.NotEqual(t, 0, len(pages), "pageChann can't be empty")
				assert.Equal(t, 0, len(posts), "postsChann must be empty")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageChan := currentTest.createPageChan(10)
			postsChan := make(chan []dto.Post, 10)

			fetchAllPosts(tt.ctx, tt.useCase, tt.jobsCounter, pageChan, postsChan)
			close(postsChan)

			tt.assertFunc(t, pageChan, postsChan)
		})
	}
}
