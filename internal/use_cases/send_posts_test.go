package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
	"api-consumer-service/internal/use_cases/mocks"
)

type sendPostsTestSuite struct{}

func (r *sendPostsTestSuite) MockDispatcherSuccess(t *testing.T) interfaces.PostsDispatcher {
	m := mocks.NewMockPostsDispatcher(t)
	m.On("Dispatch", mock.Anything, mock.Anything).Return(nil)
	return m
}

func (r *sendPostsTestSuite) MockDispatcherFailure(t *testing.T) interfaces.PostsDispatcher {
	m := mocks.NewMockPostsDispatcher(t)
	m.On("Dispatch", mock.Anything, mock.Anything).Return(errors.New("dispatch error"))

	return m
}

func TestSendPostsUseCase_Handle(t *testing.T) {
	currentTest := &sendPostsTestSuite{}

	tests := []struct {
		name       string
		ctx        context.Context
		dispatcher interfaces.PostsDispatcher
		assertFunc func(t *testing.T, ch <-chan []dto.Post, e error)
	}{
		{
			name:       "success",
			ctx:        context.Background(),
			dispatcher: currentTest.MockDispatcherSuccess(t),
			assertFunc: func(t *testing.T, ch <-chan []dto.Post, e error) {
				assert.NoError(t, e)
			},
		},
		{
			name:       "failure",
			ctx:        context.Background(),
			dispatcher: currentTest.MockDispatcherFailure(t),
			assertFunc: func(t *testing.T, ch <-chan []dto.Post, e error) {
				assert.ErrorContains(t, e, "dispatch error")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan []dto.Post)
			u := &SendPostsUseCase{
				dispatcher: tt.dispatcher,
			}
			err := u.Handle(tt.ctx, ch)

			tt.assertFunc(t, ch, err)
		})
	}
}
