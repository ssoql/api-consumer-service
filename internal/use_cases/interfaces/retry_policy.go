package interfaces

import "context"

type Retryer interface {
	Retry(ctx context.Context, operation func() error) error
}
