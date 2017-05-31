package fresh

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Route structure
type Route struct {
	path       []string
	handlers   []Handler
	params     []string
	parent     *Route
	children   []*Route
	middleware []HandlerFunc
}

type Handler struct {
	method     string
	Handler    HandlerFunc
	before     []HandlerFunc
	after      []HandlerFunc
	middelware []HandlerFunc
}

// Router structure
type Router struct {
	routes []*Route
}

// Router main function. Find the matching route and call registered handlers.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var tree func([]*Route) (bool, error)
	tree = func(routes []*Route) (bool, error) {
		for _, route := range routes {
			if strings.Join(route.path, "/") == strings.Trim(request.RequestURI, "/") {
				for _, handler := range route.handlers {
					if handler.method == request.Method {
						return true, handler.Handler(&Context{
							Request:  NewRequest(request),
							Response: NewResponse(writer),
						})
					}

				}
			}
			return tree(route.children)
		}
		return false, nil
	}
	if found, err := tree(r.routes); found && err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else if !found {
		writer.WriteHeader(http.StatusNotFound)
	}

}

// Register a route with its handlers
func (r *Router) register(method string, path string, group *Route, handlers ...HandlerFunc) error {
	var new func(*Route, string, string, ...HandlerFunc) *Route
	new = func(parentRoute *Route, method string, path string, handlers ...HandlerFunc) *Route {
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
			// handlers and middleware association
			parentRoute.addHandler(method, handlers[0], handlers[1:]...)
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
		return new(parentRoute, method, path, handlers...)
	}

	// check a group
	if group != nil {
		groupPath := strings.Join(group.path, "/")
		handlers = append(handlers[:1], append(group.middleware, handlers[1:]...)...)
		fmt.Println(handlers)
		path = groupPath + "/" + path
	}
	new(nil, method, strings.Trim(path, "/"), handlers...)
	return nil
}

// Add handlers in a route
func (route *Route) addHandler(method string, handler HandlerFunc, middleware ...HandlerFunc) {
	// If already exist an entry for the method change related handler
	changeHandler := func() bool {
		for _, h := range route.handlers {
			if h.method == method {
				h.Handler = handler
				h.middelware = append(h.middelware, middleware...)
				return false
			}
		}
		return true
	}
	if changeHandler() {
		newHandler := Handler{method: method, Handler: handler}
		newHandler.middelware = append(newHandler.middelware, middleware...)
		route.handlers = append(route.handlers, newHandler)
	}
}

// Print the list of routes
func (r *Router) printRoutes() {
	var tree func([]*Route) error
	tree = func(routes []*Route) error {
		for _, route := range routes {
			for _, handler := range route.handlers {
				log.Println(
					handler.method,
					strings.Join(route.path, "/"),
					getFuncName(handler.Handler),
					len(handler.middelware),
				)
			}
			tree(route.children)
		}
		return nil
	}
	tree(r.routes)
}

func getFuncName(f interface{}) string {
	path := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name := strings.Split(path, "/")
	return name[len(name)-1]
}
