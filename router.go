package fresh

import (
	"net/http"
)

type Route struct {

	method 	string
	name 	string
	path 	string
	handler func(Request, Response)
	before	func()
	after 	func()
}

type Router struct {
	routes []*Route
}


func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes{
		if req.RequestURI == route.path {
			route.handler(NewRequest(req), NewResponse(w))
		}
	}
	// find the right route that match current request
	// ger route and payload parameters
	// call route handler
}