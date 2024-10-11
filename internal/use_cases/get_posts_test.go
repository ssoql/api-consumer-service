package use_cases

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
	"api-consumer-service/internal/use_cases/mocks"
)

type getPostsTestSuite struct{}

func (r *getPostsTestSuite) MockApiClientSuccess(t *testing.T) interfaces.ApiConsumer {
	m := mocks.NewMockApiConsumer(t)
	m.On("DoRequest", mock.Anything, mock.Anything).Return(func(request *http.Request, posts any) error {
		tp, ok := posts.(*dto.Posts)
		assert.True(t, ok, fmt.Sprintf("incorrect type: %T", posts))

		tp.Posts = append(tp.Posts, dto.Post{
			ID:    1,
			Title: "test post",
			Body:  "test body",
		})

		return nil
	})

	return m
}

func (r *getPostsTestSuite) MockApiClientNoResults(t *testing.T) interfaces.ApiConsumer {
	m := mocks.NewMockApiConsumer(t)
	m.On("DoRequest", mock.Anything, mock.Anything).Return(func(request *http.Request, posts any) error {
		tp, ok := posts.(*dto.Posts)
		assert.True(t, ok, fmt.Sprintf("incorrect type: %T", posts))
		tp.Posts = make([]dto.Post, 0)

		return nil
	})

	return m
}

func (r *getPostsTestSuite) MockApiClientFailure(t *testing.T) interfaces.ApiConsumer {
	m := mocks.NewMockApiConsumer(t)
	m.On("DoRequest", mock.Anything, mock.Anything).Return(func(request *http.Request, total any) error {
		return errors.New("api error")
	})
	return m
}

func (r *getPostsTestSuite) MockRetryPolicy(t *testing.T) interfaces.Retryer {
	m := mocks.NewMockRetryer(t)
	m.On("Retry", mock.Anything).Return(func(operation func() error) error {
		return operation()
	})

	return m
}

func TestGetPostsUseCase_Handle(t *testing.T) {
	currentTest := &getPostsTestSuite{}

	tests := []struct {
		name          string
		ctx           context.Context
		client        interfaces.ApiConsumer
		retryStrategy interfaces.Retryer
		assertFunc    func(t *testing.T, p []dto.Post, e error)
	}{
		{
			name:          "success",
			ctx:           context.Background(),
			client:        currentTest.MockApiClientSuccess(t),
			retryStrategy: currentTest.MockRetryPolicy(t),
			assertFunc: func(t *testing.T, p []dto.Post, e error) {
				assert.NoError(t, e)
				assert.Equal(t, 1, len(p))
			},
		},
		{
			name:          "no-results",
			ctx:           context.Background(),
			client:        currentTest.MockApiClientNoResults(t),
			retryStrategy: currentTest.MockRetryPolicy(t),
			assertFunc: func(t *testing.T, p []dto.Post, e error) {
				assert.NoError(t, e)
				assert.Equal(t, 0, len(p))
			},
		},
		{
			name:          "api-error",
			ctx:           context.Background(),
			client:        currentTest.MockApiClientFailure(t),
			retryStrategy: currentTest.MockRetryPolicy(t),
			assertFunc: func(t *testing.T, p []dto.Post, e error) {
				assert.ErrorContains(t, e, "api error")
				assert.Equal(t, 0, len(p))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &GetPostsUseCase{
				client:        tt.client,
				retryStrategy: tt.retryStrategy,
			}
			got, err := u.Handle(tt.ctx, "")
			tt.assertFunc(t, got, err)
		})
	}
}
