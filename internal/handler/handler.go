package handler

import "net/http"

type authHandler interface {
	RegisterUser(http.ResponseWriter, *http.Request)
}

type Handlers struct {
	AuthHandler authHandler
}
