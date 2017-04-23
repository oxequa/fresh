package fresh

import (
	"net/http"
	"strings"
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

func (r *Router) Register(m string, p string, h func(Request, Response)) error{
	route := &Route{
		method:	m,
		path: p,
		handler: h}
	r.routes = append(r.routes, route)
	return nil
}


// Router main function. Find the matching route and call registered handlers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes{
		if req.Method == route.method && strings.TrimRight(req.RequestURI, "/") == strings.TrimRight(route.path, "/") {
			route.handler(NewRequest(req), NewResponse(w))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	// find the right route that match current request
	// ger route and payload parameters
	// call route handler
}