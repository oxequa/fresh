package fresh

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Handler structure
type (
	Handler interface {
		After(...HandlerFunc) Handler
		Before(...HandlerFunc) Handler
	}
	handler struct {
		method string
		ctrl   HandlerFunc
		before []HandlerFunc
		after  []HandlerFunc
	}
)

// Route structure
type route struct {
	path     []string
	handlers []*handler
	params   []string
	parent   *route
	children []*route
	after    []HandlerFunc
	before   []HandlerFunc
}

// Router structure
type router struct {
	routes []*route
}

// After middleware
func (h *handler) After(middleware ...HandlerFunc) Handler {
	h.after = append(h.after, middleware...)
	return h
}

// Before middleware
func (h *handler) Before(middleware ...HandlerFunc) Handler {
	h.before = append(h.before, middleware...)
	return h
}

// Router main function. Find the matching route and call registered handlers.
func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var tree func([]*route) (bool, error)
	tree = func(routes []*route) (bool, error) {
		for _, route := range routes {
			if strings.Join(route.path, "/") == strings.Trim(request.RequestURI, "/") {
				for _, handler := range route.handlers {
					if handler.method == request.Method {
						return true, handler.ctrl(&Context{
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
func (r *router) register(method string, path string, group *route, handler HandlerFunc) Handler {
	if group != nil {
		// route middleware after group middleware
		path = strings.Trim(strings.Trim(strings.Join(group.path, "/"), "/")+"/"+strings.Trim(path, "/"), "/")
		response := r.scan(nil, method, path, handler)
		if group.after != nil {
			response.After(group.after...)
		}
		if group.before != nil {
			response.Before(group.before...)
		}
		return response
	}
	return r.scan(nil, method, path, handler)
}

// Scan the routes tree
func (r *router) scan(parent *route, method string, path string, handler HandlerFunc) Handler {
	pathNodes := []string{}
	if parent != nil {
		pathNodes = strings.Split(path, "/")
		if len(pathNodes) == len(parent.path) {
			pathNodes = []string{}
		} else {
			pathNodes = pathNodes[len(parent.path):]
		}
	} else {
		pathNodes = strings.Split(path, "/")
	}
	if len(pathNodes) == 0 {
		// handlers and middleware association
		return parent.add(method, handler)
	}

	found := false
	if parent != nil {
		for _, route := range parent.children {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parent = route
				found = true
				break
			}
		}
		if found != true {
			newRoute := &route{
				path:   append(parent.path, pathNodes[0]),
				parent: parent,
			}
			parent.children = append(parent.children, newRoute)
			parent = newRoute
		}
	} else {
		for _, route := range r.routes {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parent = route
				found = true
				break
			}
		}
		if found != true {
			newRoute := &route{
				path:   []string{pathNodes[0]},
				parent: parent,
			}
			r.routes = append(r.routes, newRoute)
			parent = newRoute
		}
	}
	return r.scan(parent, method, path, handler)
}

// Add handlers to a route
func (r *route) add(method string, controller HandlerFunc, middleware ...HandlerFunc) Handler {
	// If already exist an entry for the method change related handler
	for _, h := range r.handlers {
		if h.method == method {
			h.ctrl = controller
			return h
		}
	}
	new := handler{method: method, ctrl: controller}
	r.handlers = append(r.handlers, &new)
	return &new
}

// Print the list of routes
func (r *router) printRoutes() {
	var tree func([]*route) error
	tree = func(routes []*route) error {
		for _, route := range routes {
			for _, handler := range route.handlers {
				log.Println(
					handler.method,
					strings.Join(route.path, "/"),
					getFuncName(handler.ctrl),
					len(handler.after),
					len(handler.before),
				)
			}
			tree(route.children)
		}
		return nil
	}
	tree(r.routes)
}

// Return the func name
func getFuncName(f interface{}) string {
	path := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name := strings.Split(path, "/")
	return name[len(name)-1]
}
