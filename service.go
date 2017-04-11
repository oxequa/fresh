package fresh

import(
	"net/http"
)

type service struct {
	server *http.Server
	handler *Handler
	router *Router
}