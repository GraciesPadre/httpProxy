package command

import "net/http"

type InterceptingCommand interface {
	Execute(responseWriter http.ResponseWriter, request *http.Request) (err error, handled bool)
}
