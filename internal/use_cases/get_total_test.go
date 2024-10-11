package use_cases

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
	"api-consumer-service/internal/use_cases/mocks"
)

type getTotalTestSuite struct{}

func (r *getTotalTestSuite) MockApiClientSuccess(t *testing.T) interfaces.ApiConsumer {
	m := mocks.NewMockApiConsumer(t)
	m.On("DoRequest", mock.Anything, mock.Anything).Return(func(request *http.Request, total any) error {
		tp, ok := total.(*dto.ResponseTotal)
		assert.True(t, ok)
		tp.Total = 1

		return nil
	})

	return m
}

func (r *getTotalTestSuite) MockApiClientFailure(t *testing.T) interfaces.ApiConsumer {
	m := mocks.NewMockApiConsumer(t)
	m.On("DoRequest", mock.Anything, mock.Anything).Return(func(request *http.Request, total any) error {
		return errors.New("api error")
	})
	return m
}

func (r *getTotalTestSuite) MockRetryPolicy(t *testing.T) interfaces.Retryer {
	m := mocks.NewMockRetryer(t)
	m.On("Retry", mock.Anything).Return(func(operation func() error) error {
		return operation()
	})

	return m
}

func TestGetTotalUseCase_GetTotal(t *testing.T) {
	currentTest := &getTotalTestSuite{}

	tests := []struct {
		name          string
		ctx           context.Context
		client        interfaces.ApiConsumer
		retryStrategy interfaces.Retryer
		assertFunc    func(t *testing.T, total int, e error)
	}{
		{
			name:          "success",
			ctx:           context.Background(),
			client:        currentTest.MockApiClientSuccess(t),
			retryStrategy: currentTest.MockRetryPolicy(t),
			assertFunc: func(t *testing.T, total int, e error) {
				assert.NoError(t, e)
				assert.Equal(t, 1, total)
			},
		},
		{
			name:          "api-error",
			ctx:           context.Background(),
			client:        currentTest.MockApiClientFailure(t),
			retryStrategy: currentTest.MockRetryPolicy(t),
			assertFunc: func(t *testing.T, total int, e error) {
				assert.ErrorContains(t, e, "api error")
				assert.Equal(t, 0, total)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &GetTotalUseCase{
				client:        tt.client,
				retryStrategy: tt.retryStrategy,
			}
			got, err := u.GetTotal(tt.ctx, "")
			tt.assertFunc(t, got, err)
		})
	}
}
