package infrastructure

import (
	"context"
	"errors"
	"log"
	"time"
)

// RetryWithBackoff implements a retry policy with exponential backoff
type RetryWithBackoff struct {
	retries int
}

// NewExponentialBackoff creates a new instance of RetryWithBackoff
func NewExponentialBackoff(retry int) *RetryWithBackoff {
	return &RetryWithBackoff{retry}
}

// Retry retries the operation with exponential backoff on failure
func (r *RetryWithBackoff) Retry(operation func() error) error {
	var err error

	if r.retries < 1 {
		return errors.New("number of retries should be higher or equal 1")
	}

	for attempt := 1; attempt <= r.retries; attempt++ {
		if err = operation(); err == nil {
			return nil
		}
		log.Printf("Retry attempt %d failed: %v", attempt, err)

		if errors.Is(err, context.Canceled) {
			return err
		}

		time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
	}

	return err
}
