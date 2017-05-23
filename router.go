package fresh

import (
	"log"
	"net/http"
	"strings"
)

// Route structure
type Route struct {
	path     []string
	handlers []Handler
	params   []string
	parent   *Route
	children []*Route
}

type Handler struct {
	method     string
	Handler    HandlerFunc
	before     []MiddlewareFunc
	after      []MiddlewareFunc
	middelware []MiddlewareFunc
}

// Router structure
type Router struct {
	routes []*Route
}

func (r *Router) Register(method string, path string, handler HandlerFunc) error {
	r.newRoute(nil, method, strings.Trim(path, "/"), handler)
	return nil
}

func (r *Router) newRoute(parentRoute *Route, method string, path string, handler HandlerFunc) *Route {
	pathNodes := []string{}

	if parentRoute != nil {
		pathNodes = strings.Split(path, "/")
		if len(pathNodes) == len(parentRoute.path) {
			pathNodes = []string{}
		} else {
			pathNodes = pathNodes[len(parentRoute.path):]
		}
	} else {
		pathNodes = strings.Split(path, "/")
	}
	if len(pathNodes) == 0 {
		parentRoute.addHandler(method, handler)
		return parentRoute
	}
	found := false
	if parentRoute != nil {
		for _, route := range parentRoute.children {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parentRoute = route
				found = true
				break
			}
		}
		if found != true {
			newRoute := &Route{
				path:   append(parentRoute.path, pathNodes[0]),
				parent: parentRoute,
			}
			parentRoute.children = append(parentRoute.children, newRoute)
			parentRoute = newRoute

		}
	} else {
		for _, route := range r.routes {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parentRoute = route
				found = true
				break
			}
		}
		if found != true {
			newRoute := &Route{
				path:   []string{pathNodes[0]},
				parent: parentRoute,
			}
			r.routes = append(r.routes, newRoute)
			parentRoute = newRoute
		}
	}

	return r.newRoute(parentRoute, method, path, handler)
}

// Router main function. Find the matching route and call registered handlers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Add handler in a route
func (route *Route) addHandler(method string, handler HandlerFunc) {
	// if already exist an entry for the method change related handler
	changeHandler := func() bool {
		for _, h := range route.handlers {
			if h.method == method {
				h.Handler = handler
				return false
			}
		}
		return true
	}
	if changeHandler() {
		newHandler := Handler{method: method, Handler: handler}
		route.handlers = append(route.handlers, newHandler)
	}
}

func (r *Router) PrintRoutes() {
	var tree func([]*Route) error
	tree = func(routes []*Route) error {
		for _, route := range routes {
			for _, handler := range route.handlers {
				log.Println(handler.method + " - " + strings.Join(route.path, "/"))
			}
			return tree(route.children)
		}
		return nil
	}
	tree(r.routes)
}
