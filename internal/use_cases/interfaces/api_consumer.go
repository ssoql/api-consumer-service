package interfaces

import "net/http"

type ApiConsumer interface {
	DoRequest(request *http.Request, results any) error
}
