package interfaces

type Retryer interface {
	Retry(operation func() error) error
}
