package use_cases

import (
	"context"
	"net/http"

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
	var total int
	type responseTotal struct {
		Total int `json:"total"`
	}
	respTotal := responseTotal{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return total, err
	}

	err = u.client.DoRequest(req, &respTotal)
	if err != nil {
		return total, err
	}

	return respTotal.Total, nil
}
