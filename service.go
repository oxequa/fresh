package fresh

import (
	"net/http"
)

type service struct {
	server  *http.Server
	router  *Router
}
