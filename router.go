package fresh

import (
	"net/http"
)

type Route struct {
	method string
	name string
	path string
	handler http.Handler
	before func()
	after func()
}

type Router struct {
	routes []*Route
}

func (r *Router) add(m string, p string, h func(), b func(), a func()) {

}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}