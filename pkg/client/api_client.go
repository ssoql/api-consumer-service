package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApiClient struct {
	retry  int
	client *http.Client
}

func NewApiClient(timeout int) *ApiClient {
	return &ApiClient{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *ApiClient) DoRequest(request *http.Request, results any) error {
	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(results)
	resp.Body.Close()

	if err != nil {
		return err
	}

	return nil
}
