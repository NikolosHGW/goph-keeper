package handler

import "net/http"

type registerHandler interface {
	RegisterUser(http.ResponseWriter, *http.Request)
}

type Handlers struct {
	RegisterHandler registerHandler
}
