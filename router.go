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
	path     	[]string
	handlers 	[]*handler
	parent   	 *route
	children 	[]*route
	after    	[]HandlerFunc
	before   	[]HandlerFunc
}

// Router struct
type router struct {
	routes []*route
	parameters  map[string] string
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
	h := handler{method: method, ctrl: controller}
	r.handlers = append(r.handlers, &h)
	return &h
}

// Scan the routes tree
func (r *router) scan(parent *route, method string, path string, handler HandlerFunc) Handler {
	pathNodes := []string{}
	pathNodes = strings.Split(strings.Trim(path, "/"), "/")

	if len(pathNodes) == 1 && pathNodes[0] ==  ""{
		return nil
	}
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

	// TODO: refactor
	found := false
	if parent != nil {
		for i, route := range parent.children {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parent = route
				found = true
				break
			}
			if i == len(parent.children) - 1 && isURLParameter(pathNodes[0]) && isURLParameter(parent.children[i].path[len(parent.children[i].path) -1]){
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
			if isURLParameter(pathNodes[0]) == false && len(parent.children) > 0 && isURLParameter(parent.children[len(parent.children) - 1].path[len(parent.children[len(parent.children) - 1].path) - 1]) == true{
				if len(parent.children) > 1 {
					parent.children = append(parent.children[:len(parent.children) - 1], newRoute, parent.children[len(parent.children) - 1])
				} else {

					parent.children = append([] *route{newRoute}, parent.children...)
				}
			} else {
				parent.children = append(parent.children, newRoute)
			}
			parent = newRoute
		}
	} else {
		for i, route := range r.routes {
			if route.path[len(route.path)-1] == pathNodes[0] {
				parent = route
				found = true
				break
			}
			if i == len(r.routes) - 1 && isURLParameter(pathNodes[0]) && isURLParameter(r.routes[i].path[len(r.routes[i].path) -1]){
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
			if isURLParameter(pathNodes[0]) == false && len(r.routes) > 0 && isURLParameter(r.routes[len(r.routes) - 1].path[len(r.routes[len(r.routes) - 1].path) - 1]) == true{
				if len(r.routes) > 1 {
					r.routes = append(r.routes[:len(r.routes) - 1], newRoute, r.routes[len(r.routes) - 1])
				} else {

					r.routes = append([] *route{newRoute}, r.routes...)
				}
			} else {
				r.routes = append(r.routes, newRoute)
			}
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

// Process a request
func (r *router) process(handler *handler, context *context, ws bool) (err error){
	if err := handler.middleware(context, handler.before...); err != nil {
		return err
	}
	if ws{
		//TODO
	}else {
		err = handler.ctrl(context)
	}
	if err != nil {
		return err
	}
	// after middleware
	if err := handler.middleware(context, handler.after...); err != nil {
		return err
	}
	// write response
	context.response.write()
	return nil
}

// Tree return a matching route if exist
func (r *router) tree(routes []*route, index int, response http.ResponseWriter, request *http.Request) (bool, error){
	context := new(context)
	context.init(request, response)
	r.parameters = make(map[string] string)
	for i, route := range routes {
		if index > len(strings.Split(strings.Trim(request.URL.Path, "/"), "/")) - 1 {
			return false, nil
		}
		urlPath := strings.Split(strings.Trim(request.URL.Path, "/"), "/")[index]
		if route.path[len(route.path) - 1] == urlPath {
			if strings.Join(route.path, "/") == strings.Trim(request.URL.Path, "/") {
				for _, handler := range route.handlers {
					switch {
					case context.request.IsWS():
						// web socket call
						return true, r.process(handler,context, true)
					case handler.method == request.Method:
						// http call
						return true, r.process(handler,context, false)
					default:
						break
					}
				}
			}
			return r.tree(route.children, index+1, response, request)
		} else if i == len(routes) - 1 {
			parameterName := routes[i].path[len(routes[i].path) - 1]
			if isURLParameter(parameterName) == true {
				r.parameters[strings.Trim(strings.Trim(parameterName, "{"), "}")] = urlPath
				if index == len(strings.Split(strings.Trim(request.URL.Path, "/"), "/")) - 1{
					for _, handler := range route.handlers {
						switch {
						case context.request.IsWS():
							// web socket call
							return true, r.process(handler,context, true)
						case handler.method == request.Method:
							// http call
							return true, r.process(handler,context, false)
						default:
							break
						}
					}
				}
				return r.tree(route.children, index+1,response,request)
			}

		}
	}
	return false, nil
}

// After middleware for a single route
func (h *handler) After(middleware ...HandlerFunc) Handler {
	if middleware != nil {
		h.after = append(h.after, middleware...)
	}
	return h
}

// Before middleware for a single route
func (h *handler) Before(middleware ...HandlerFunc) Handler {
	if middleware != nil {
		h.before = append(h.before, middleware...)
	}
	return h
}

// After middleware for a resource
func (r *resource) After(middleware ...HandlerFunc) Resource {
	if middleware != nil {
		for _, route := range r.rest {
			route.After(middleware...)
		}
	}
	return r
}

// Before middleware for a resource
func (r *resource) Before(middleware ...HandlerFunc) Resource {
	if middleware != nil {
		for _, route := range r.rest {
			route.Before(middleware...)
		}
	}
	return r
}

// Router main function. Find the matching route and call registered handlers.
func (r *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if found, err := r.tree(r.routes, 0, response, request); found && err != nil {
		// Handle internal server error
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	} else if !found {
		response.WriteHeader(http.StatusNotFound)
	}
}

func isURLParameter(value string) bool {
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}"){
		return true
	}
	return false
}