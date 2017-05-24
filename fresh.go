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
		Get(string, HandlerFunc) error
		Post(string, HandlerFunc) error
		Put(string, HandlerFunc) error
		Trace(string, HandlerFunc) error
		Patch(string, HandlerFunc) error
		Delete(string, HandlerFunc) error
		Options(string, HandlerFunc) error
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
		f.router.PrintRoutes()
		f.server.Serve(listener)
	}()
	<-shutdown
	log.Println("Server shutting down")
	ctx, _ := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	return nil
}

// Routes group
func (f *fresh) Group(path string, handlers ...HandlerFunc) Fresh {
	f.group = &Route{}
	f.group.path = strings.Split(path, "/")
	f.group.middleware = append(f.group.middleware, handlers...)
	return f
}

// Register for GET APIs
func (f *fresh) Get(path string, handler HandlerFunc) error {
	return f.router.Register("GET", path, handler)
}

// Register for POST APIs
func (f *fresh) Post(path string, handler HandlerFunc) error {
	return f.router.Register("POST", path, handler)
}

// Register for PUT APIs
func (f *fresh) Put(path string, handler HandlerFunc) error {
	return f.router.Register("PUT", path, handler)
}

// Register for PATCH APIs
func (f *fresh) Patch(path string, handler HandlerFunc) error {
	return f.router.Register("PATCH", path, handler)
}

// Register for DELETE APIs
func (f *fresh) Delete(path string, handler HandlerFunc) error {
	return f.router.Register("DELETE", path, handler)
}

// Register for OPTIONS APIs
func (f *fresh) Options(path string, handler HandlerFunc) error {
	return f.router.Register("OPTIONS", path, handler)
}

// Register for TRACE APIs
func (f *fresh) Trace(path string, handler HandlerFunc) error {
	return f.router.Register("TRACE", path, handler)
}
