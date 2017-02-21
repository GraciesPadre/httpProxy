package command

import "net/http"

type Command interface {
	Execute(responseWriter http.ResponseWriter, request *http.Request) (err error, handled bool)
}
