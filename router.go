package fresh

import (
	"net/http"
)


// Route structure
type Route struct {
	method 	string
	name 	string
	path 	string
	handler func(Request, Response)
	before	func()
	after 	func()
}


// Router structure
type Router struct {
	routes []*Route
}


// Router main function. Find the matching route and call registered handlers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes{
		if req.Method == route.method && req.RequestURI == route.path {
			route.handler(NewRequest(req), NewResponse(w))
		}
	}

	w.WriteHeader(http.StatusNotFound)
	// find the right route that match current request
	// ger route and payload parameters
	// call route handler
}