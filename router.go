package fresh

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	path      string
	handlers  []*handler
	parent    *route
	children  []*route
	after     []HandlerFunc
	before    []HandlerFunc
	parameter bool
}

// Router struct
type router struct {
	*fresh
	route  *route
	static map[string]string
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

// isURLParameter check if given string is a param
func isURLParameter(value string) bool {
	if strings.HasPrefix(value, ":") {
		return true
	}
	return false
}

// TODO remove
func getFuncName(f interface{}) string {
	path := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name := strings.Split(path, "/")
	return name[len(name)-1]
}

func (r *route) getHandler(method string) *handler {
	for _, h := range r.handlers {
		if h.method == method {
			return h
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

// Process a request
func (r *router) process(handler *handler, response http.ResponseWriter, request *http.Request, context *context) (err error) {
	context.init(request, response)
	if err = handler.middleware(context, handler.before...); err != nil {
		return err
	}
	// route controller
	err = handler.ctrl(context)
	if err != nil {
		return err
	}
	// after middleware
	if err = handler.middleware(context, handler.after...); err != nil {
		return err
	}

	for _, ch := range r.Config.handlers {
		err := ch(context)
		if err != nil {
			return err
		}
	}
	// write response
	context.response.write()
	return
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

// Print the list of routes
func (r *router) printRoutes() {
	var tree func(routes []*route, parentPath string) error
	tree = func(routes []*route, parentPath string) error {
		for _, route := range routes {
			separator := ""
			if strings.HasSuffix(parentPath, "/") == false {
				separator = "/"
			}
			currentPath := parentPath + separator + route.path

			for _, handler := range route.handlers {
				log.Println(
					handler.method,
					currentPath,
					getFuncName(handler.ctrl),
					len(handler.before),
					len(handler.after),
				)
			}
			tree(route.children, currentPath)
		}
		return nil
	}
	tree([]*route{r.route}, "")
}

// Register static routes for assets
func (r *router) addStatic(static map[string]string) Handler {
	for k, v := range static {
		r.static[k] = v
	}
	return nil
}

// Router main function. Find the matching route and call registered handlers.
func (r *router) serveStatic(response http.ResponseWriter, request *http.Request) {
	for publicPath, staticPath := range r.static {
		path := strings.Replace(strings.Trim(request.URL.Path, "/"), publicPath, staticPath, 1)
		path, _ = filepath.Abs(path)
		f, err := os.Stat(path)
		if err == nil && !f.IsDir() {
			http.ServeFile(response, request, path)
			return
		} else if f.IsDir() {
			for _, testDefaultFile := range r.Config.Default {
				filePath := filepath.Join(path, testDefaultFile)
				if f, err := os.Stat(filePath); err == nil && !f.IsDir() {
					http.ServeFile(response, request, filePath)
					return
				}
			}

		}
	}
	http.NotFound(response, request)
}

// Register a route with its handlers
func (r *router) register(method string, path string, handler HandlerFunc) Handler {
	splittedPath := strings.Split(strings.Trim(path, "/"), "/")
	route := r.scanTree(r.route, splittedPath, nil, true)
	return route.add(method, handler)
}

// Scan the tree to find the matching route (if save create all needed routes)
func (r *router) scanTree(parent *route, path []string, context *context, save bool) *route {
	if len(path) > 0 {
		for _, route := range parent.children {
			if route.path == path[0] {
				return r.scanTree(route, path[1:], context, save)
			}
			if !save && route.parameter {
				context.parameters[route.path[1:]] = path[0]
				return r.scanTree(route, path[1:], context, save)
			}
		}
		if !save {
			if parent.children[len(parent.children)-1].parameter {
				return parent.children[len(parent.children)-1]
			} else {
				return nil
			}
		}
		newRoute := &route{path: path[0], parent: parent}
		switch {
		case isURLParameter(path[0]):
			newRoute.parameter = true
			parent.children = append(parent.children, newRoute)
		default:
			parent.children = append([]*route{newRoute}, parent.children...)
		}
		return r.scanTree(newRoute, path[1:], context, save)
	}
	return parent
}

// Router main function. Find the matching route and call registered handlers.
func (r *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context := &context{}
	context.parameters = make(map[string]string)
	splittedPath := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	if route := r.scanTree(r.route, splittedPath, context, false); route != nil {
		if r.Options && request.Method == "OPTIONS" {
			h := &handler{
				ctrl: func(c Context) error {
					return c.Response().Code(http.StatusOK)
				},
			}
			r.process(h, response, request, context)
			return
		}
		if routeHandler := route.getHandler(request.Method); routeHandler != nil {
			err := r.process(routeHandler, response, request, context)
			if err != nil {
				context.Response().writeErr(err)
			}
			return
		}
	}
	r.serveStatic(response, request)
}
