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
)

// Main Fresh structure
type (
	Fresh interface {
		Run() error
		Group(string, ...HandlerFunc) Fresh
		Resource(string, ...HandlerFunc) error
		Get(string, ...HandlerFunc) error
		Post(string, ...HandlerFunc) error
		Put(string, ...HandlerFunc) error
		Trace(string, ...HandlerFunc) error
		Patch(string, ...HandlerFunc) error
		Delete(string, ...HandlerFunc) error
		Options(string, ...HandlerFunc) error
	}

	fresh struct {
		group  *Route
		config *config
		router *Router
		server *http.Server
	}

	Context struct {
		Request  Request
		Response Response
	}

	HandlerFunc func(*Context) error
)

// Initialize main Fresh structure
func New() Fresh {
	fresh := fresh{
		config: new(config),
		server: new(http.Server),
		router: new(Router),
	}
	wd, _ := os.Getwd()
	if fresh.config.read(wd) != nil {
		// random ip and port
		rand.Seed(time.Now().Unix())
		fresh.config.Host = "localhost"
		fresh.config.Port = rand.Intn(9999-1111) + 1111
	}
	return &fresh
}

// Load all servers configurations and start them
func (f *fresh) Run() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.config.Host+":"+strconv.Itoa(f.config.Port))
	if err != nil {
		return err
	}
	go func() {
		log.Println("Server started on ", f.config.Host, ":", f.config.Port)
		f.server.Handler = f.router
		f.router.printRoutes()
		f.server.Serve(listener)
	}()
	<-shutdown
	log.Println("Server shutting down")
	ctx, _ := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	return nil
}

// Register a group
func (f fresh) Group(path string, middleware ...HandlerFunc) Fresh {
	f.group = &Route{}
	f.group.path = strings.Split(path, "/")
	f.group.middleware = append(f.group.middleware, middleware...)
	return &f
}

// Register a resource (get, post, put, delete)
func (f *fresh) Resource(path string, handlers ...HandlerFunc) (err error) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	path = strings.Trim(path, "/")
	name := path
	if strings.LastIndex(path, "/") != -1 {
		name = string(path[strings.LastIndex(path, "/")+1:])
	}
	name = "{" + name + "}"
	for _, method := range methods {
		switch method {
		case "GET":
			err = f.router.register(method, path, f.group, handlers[0])
		case "POST":
			err = f.router.register(method, path+"/"+name, f.group, handlers[1])
		case "PUT", "PATCH":
			err = f.router.register(method, path+"/"+name, f.group, handlers[2])
		case "DELETE":
			err = f.router.register(method, path+"/"+name, f.group, handlers[3])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Register for GET APIs
func (f *fresh) Get(path string, handlers ...HandlerFunc) error {
	return f.router.register("GET", path, f.group, handlers...)
}

// Register for POST APIs
func (f *fresh) Post(path string, handlers ...HandlerFunc) error {
	return f.router.register("POST", path, f.group, handlers...)
}

// Register for PUT APIs
func (f *fresh) Put(path string, handlers ...HandlerFunc) error {
	return f.router.register("PUT", path, f.group, handlers...)
}

// Register for PATCH APIs
func (f *fresh) Patch(path string, handlers ...HandlerFunc) error {
	return f.router.register("PATCH", path, f.group, handlers...)
}

// Register for DELETE APIs
func (f *fresh) Delete(path string, handlers ...HandlerFunc) error {
	return f.router.register("DELETE", path, f.group, handlers...)
}

// Register for OPTIONS APIs
func (f *fresh) Options(path string, handlers ...HandlerFunc) error {
	return f.router.register("OPTIONS", path, f.group, handlers...)
}

// Register for TRACE APIs
func (f *fresh) Trace(path string, handlers ...HandlerFunc) error {
	return f.router.register("TRACE", path, f.group, handlers...)
}
