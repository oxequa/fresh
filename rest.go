package fresh

import (
	"golang.org/x/net/websocket"
	"strings"
)

type Rest interface {
	STATIC(map[string]string)
	WS(string, HandlerFunc) Handler
	GET(string, HandlerFunc) Handler
	POST(string, HandlerFunc) Handler
	PUT(string, HandlerFunc) Handler
	TRACE(string, HandlerFunc) Handler
	PATCH(string, HandlerFunc) Handler
	DELETE(string, HandlerFunc) Handler
	OPTIONS(string, HandlerFunc) Handler
	CRUD(string, ...HandlerFunc) Resource
}

// WS api registration
func (f *fresh) WS(path string, handler HandlerFunc) Handler {
	h := func(c Context) (err error) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			c.Request().SetWS(ws)
			err = handler(c)
		}).ServeHTTP(c.Response().Get(), c.Request().Get())
		return err
	}
	return f.router.addRoute("GET", path, h)
}

// Register a resource (get, post, put, delete)
func (f *fresh) CRUD(path string, h ...HandlerFunc) Resource {
	res := resource{
		methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}
	path = strings.Trim(path, "/")
	name := "{" + path + "}"
	if strings.LastIndex(path, "/") != -1 {
		name = string("{" + path[strings.LastIndex(path, "/")+1:] + "}")
	}
	for _, method := range res.methods {
		switch method {
		case "GET":
			res.rest = append(res.rest, f.router.addRoute(method, path, h[0]))
		case "POST":
			res.rest = append(res.rest, f.router.addRoute(method, path+"/"+name, h[1]))
		case "PUT", "PATCH":
			res.rest = append(res.rest, f.router.addRoute(method, path+"/"+name, h[2]))
		case "DELETE":
			res.rest = append(res.rest, f.router.addRoute(method, path+"/"+name, h[3]))
		}
	}
	return &res
}

// GET api registration
func (f *fresh) GET(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("GET", path, handler)
}

// PUT api registration
func (f *fresh) PUT(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("PUT", path, handler)
}

// POST api registration
func (f *fresh) POST(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("POST", path, handler)
}

// TRACE api registration
func (f *fresh) TRACE(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("TRACE", path, handler)
}

// PATCH api registration
func (f *fresh) PATCH(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("PATCH", path, handler)
}

// DELETE api registration
func (f *fresh) DELETE(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("DELETE", path, handler)
}

// OPTIONS api registration
func (f *fresh) OPTIONS(path string, handler HandlerFunc) Handler {
	return f.router.addRoute("OPTIONS", path, handler)
}

// ASSETS serve a list of static files. Array of files or directories TODO write logic
func (f *fresh) STATIC(static map[string]string) {
	f.router.addStatic(static)
}
