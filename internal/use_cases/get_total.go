package use_cases

import (
	"context"
	"net/http"

	"api-consumer-service/internal/dto"
	"api-consumer-service/internal/use_cases/interfaces"
)

type GetTotalUseCase struct {
	client        interfaces.ApiConsumer
	retryStrategy interfaces.Retryer
}

func NewGetTotalUseCase(client interfaces.ApiConsumer) *GetTotalUseCase {
	return &GetTotalUseCase{
		client: client,
	}
}

func (u *GetTotalUseCase) GetTotal(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	respTotal := dto.ResponseTotal{}
	callback := func() error {
		return u.client.DoRequest(req, &respTotal)
	}

	err = u.retryStrategy.Retry(callback)
	if err != nil {
		return 0, err
	}

	return respTotal.Total, nil
}
