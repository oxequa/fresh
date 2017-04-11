package fresh

import (
	"net/http"
)

type Route struct {
	name    string
	path    string
	handler http.Handler
}

type Router struct {
	routes []*Route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
