package fresh

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Handler struct
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

// Route struct
type route struct {
	path     []string
	handlers []*handler
	params   []string
	parent   *route
	children []*route
	after    []HandlerFunc
	before   []HandlerFunc
}

// Router struct
type router struct {
	routes []*route
}

// Resource struct
type (
	Resource interface {
		After(...HandlerFunc) Resource
		Before(...HandlerFunc) Resource
	}
	resource struct {
		methods []string
		rest    []Handler
	}
)

// Return the func name
func getFuncName(f interface{}) string {
	path := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name := strings.Split(path, "/")
	return name[len(name)-1]
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
					len(handler.before),
					len(handler.after),
				)
			}
			tree(route.children)
		}
		return nil
	}
	tree(r.routes)
}

// Run a middleware
func (h *handler) middleware(c Context, handlers ...HandlerFunc) error {
	for _, f := range handlers {
		if f != nil {
			if err := f(c); err != nil {
				return err
			}
		}
	}
	return nil
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

// Scan the routes tree
func (r *router) scan(parent *route, method string, path string, handler HandlerFunc) Handler {
	pathNodes := []string{}
	pathNodes = strings.Split(path, "/")
	for index, str := range pathNodes {
		if len(str) == 0 {
			 pathNodes = append(pathNodes[:index], pathNodes[index + 1:]...)
		}
	}
	if parent != nil {
		if len(pathNodes) == len(parent.path) {
			pathNodes = []string{}
		} else {
			pathNodes = pathNodes[len(parent.path):]
		}
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

// After middleware for a single route
func (h *handler) After(middleware ...HandlerFunc) Handler {
	h.after = append(h.after, middleware...)
	return h
}

// Before middleware for a single route
func (h *handler) Before(middleware ...HandlerFunc) Handler {
	h.before = append(h.before, middleware...)
	return h
}

// After middleware for a resource
func (r *resource) After(middleware ...HandlerFunc) Resource {
	for _, route := range r.rest {
		route.After(middleware...)
	}
	return r
}

// Before middleware for a resource
func (r *resource) Before(middleware ...HandlerFunc) Resource {
	for _, route := range r.rest {
		route.Before(middleware...)
	}
	return r
}

// Router main function. Find the matching route and call registered handlers.
func (r *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var tree func([]*route, int) (bool, error)
	tree = func(routes []*route, index int) (bool, error) {
		for _, route := range routes {
			if route.path[len(route.path) - 1] == strings.Split(strings.Trim(request.RequestURI, "/"), "/")[index] {
				if strings.Join(route.path, "/") == strings.Trim(request.RequestURI, "/") {
					for _, handler := range route.handlers {
						if handler.method == request.Method {
							context := context{}
							context.new(request, response)

							// before middleware
							if err := handler.middleware(&context, handler.before...); err != nil {
								return true, err
							}
							reply := handler.ctrl(&context)
							if reply != nil {
								return true, reply
							}
							// after middleware
							if err := handler.middleware(&context, handler.after...); err != nil {
								return true, err
							}
							// write response
							context.response.write()
							return true, nil
						}
					}
				}
				tree(route.children, index+1)
			}
		}
		// Get path parameters
		if index == len(strings.Split(strings.Trim(request.RequestURI, "/"), "/")) - 1{
			parameterName := routes[len(routes) - 1].path[len(routes[len(routes) - 1].path) - 1]
			if strings.HasPrefix(parameterName, "{") &&
				strings.HasSuffix(parameterName, "}"){
				log.Println(parameterName[1:len(parameterName) - 1])
				log.Println(strings.Split(strings.Trim(request.RequestURI, "/"), "/")[len(strings.Split(strings.Trim(request.RequestURI, "/"), "/")) - 1])
			}
		}
		return false, nil
	}
	if found, err := tree(r.routes, 0); found && err != nil {
		// Handle internal server error
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else if !found {
		response.WriteHeader(http.StatusNotFound)
	}

}