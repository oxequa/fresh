package fresh

import (
	"net/http"
	"strings"
	"regexp"
)


// Route structure
type Route struct {
	method 	string
	name 	string
	path 	string
	handler func(Request, Response)
	before	func()
	after 	func()
	params	[]string
}


// Router structure
type Router struct {
	routes []*Route
}

func (r *Router) Register(m string, p string, h func(Request, Response)) error{
	params := []string{}
	for {
		s := strings.Index(p, "{")
		e := strings.Index(p, "}")
		if s == -1 {
			break
		}
		if e == -1 {
			panic("URL template error: missing brackets.")
		}
		params = append(params, p[s+1:e])
		p = strings.Replace(p, p[s:e+1], ".([^/W]+)", 1)

	}
	route := &Route{
		method:	m,
		path: p,
		handler: h,
		params: params}
	r.routes = append(r.routes, route)
	return nil
}


// Router main function. Find the matching route and call registered handlers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathToMatch := strings.TrimRight(req.RequestURI, "/")
	if strings.LastIndex(pathToMatch, "?") != -1{
		pathToMatch = pathToMatch[:strings.LastIndex(pathToMatch, "?")]
	}
	for _, route := range r.routes{
		if req.Method == route.method {
			match, _ := regexp.MatchString("^" + strings.TrimRight(route.path, "/") +  "$", pathToMatch)

			if match{
				route.handler(NewRequest(req), NewResponse(w))
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	// find the right route that match current request
	// ger route and payload parameters
	// call route handler
}