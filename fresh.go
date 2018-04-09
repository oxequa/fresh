package fresh

import (
	httpContext "context"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/websocket"
)

// Main Fresh structure
type (
	fresh struct {
		*Config
		*http.Server
		router *router
	}

	context struct {
		request    request
		response   response
		parameters map[string]string
	}

	Rest interface {
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

	Fresh interface {
		Rest
		Run() error
		Shutdown() error
		Group(string) Group
	}

	Context interface {
		Request() Request
		Response() Response
		Writer(http.ResponseWriter)
	}

	HandlerFunc func(Context) error
)

// Initialize main Fresh structure
func New() Fresh {
	fresh := fresh{}
	// Fresh config
	fresh.Config = &Config{}
	fresh.Config.fresh = &fresh
	fresh.Config.Init()
	// Server setting
	fresh.Server = new(http.Server)
	// Fresh router
	fresh.router = &router{&fresh, &route{}, make(map[string]string)}

	wd, _ := os.Getwd()
	if fresh.Config.read(wd) != nil {
		// random port
		rand.Seed(time.Now().Unix())
		fresh.Config.Host = "127.0.0.1"
		fresh.Config.Port = rand.Intn(9999-1111) + 1111
	}
	return &fresh
}

// Run HTTP server
func (f *fresh) Run() error {
	shutdown := make(chan os.Signal)
	port := strconv.Itoa(f.Port)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.Host+":"+port)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go func() {
		log.Println("Server listen on", f.Host+":"+port)
		f.Server.Handler = f.router
		f.router.printRoutes()
		// check for tsl before serve
		if f.TSL != nil {
			f.tsl()
		}
		f.Server.Serve(listener)
	}()
	<-shutdown
	f.Shutdown()
	return nil
}

// Shutdown server
func (f *fresh) Shutdown() error {
	ctx, cancel := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.Server.Shutdown(ctx)
	cancel()
	log.Println("Server shutdown")
	return nil
}

// Group registration
func (f *fresh) Group(path string) Group {
	g := group{
		parent: f,
		route: &route{
			path: path,
		},
	}
	return &g
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
	return f.router.register("GET", path, h)
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
			res.rest = append(res.rest, f.router.register(method, path, h[0]))
		case "POST":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[1]))
		case "PUT", "PATCH":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[2]))
		case "DELETE":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[3]))
		}
	}
	return &res
}

// GET api registration
func (f *fresh) GET(path string, handler HandlerFunc) Handler {
	return f.router.register("GET", path, handler)
}

// PUT api registration
func (f *fresh) PUT(path string, handler HandlerFunc) Handler {
	return f.router.register("PUT", path, handler)
}

// POST api registration
func (f *fresh) POST(path string, handler HandlerFunc) Handler {
	return f.router.register("POST", path, handler)
}

// TRACE api registration
func (f *fresh) TRACE(path string, handler HandlerFunc) Handler {
	return f.router.register("TRACE", path, handler)
}

// PATCH api registration
func (f *fresh) PATCH(path string, handler HandlerFunc) Handler {
	return f.router.register("PATCH", path, handler)
}

// DELETE api registration
func (f *fresh) DELETE(path string, handler HandlerFunc) Handler {
	return f.router.register("DELETE", path, handler)
}

// OPTIONS api registration
func (f *fresh) OPTIONS(path string, handler HandlerFunc) Handler {
	return f.router.register("OPTIONS", path, handler)
}

// ASSETS serve a list of static files. Array of files or directories TODO write logic
func (f *fresh) STATIC(static map[string]string) {
	f.router.addStatic(static)
}

// Return context request
func (c *context) Request() Request {
	return &c.request
}

// Return context response
func (c *context) Response() Response {
	return &c.response
}

// Overwrite http writer
func (c *context) Writer(w http.ResponseWriter) {
	c.response.w = w
}

// Init set context request and response
func (c *context) init(r *http.Request, w http.ResponseWriter) {
	c.response = response{w: w, r: r}
	c.request = request{r: r}
	c.request.setRouteParam(c.parameters)
}
