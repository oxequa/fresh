package fresh

import (
	"context"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Main Fresh structure
type (
	Fresh interface {
		Run() error
		Get(string, HandlerFunc) error
		Post(string, HandlerFunc) error
		Put(string, HandlerFunc) error
		Trace(string, HandlerFunc) error
		Patch(string, HandlerFunc) error
		Delete(string, HandlerFunc) error
		Options(string, HandlerFunc) error
	}
	fresh struct {
		config *config
		router *Router
		server *http.Server
	}

	MiddlewareFunc func(context.Context) error

	HandlerFunc func(Request, Response) HTTPError

	HTTPError struct {
		Code int
		Body interface{}
	}
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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	return nil
}

// Register for GET APIs
func (f *fresh) Get(p string, h HandlerFunc) error {
	return f.router.Register("GET", p, h)
}

// Register for POST APIs
func (f *fresh) Post(p string, h HandlerFunc) error {
	return f.router.Register("POST", p, h)
}

// Register for PUT APIs
func (f *fresh) Put(p string, h HandlerFunc) error {
	return f.router.Register("PUT", p, h)
}

// Register for PATCH APIs
func (f *fresh) Patch(p string, h HandlerFunc) error {
	return f.router.Register("PATCH", p, h)
}

// Register for DELETE APIs
func (f *fresh) Delete(p string, h HandlerFunc) error {
	return f.router.Register("DELETE", p, h)
}

// Register for OPTIONS APIs
func (f *fresh) Options(p string, h HandlerFunc) error {
	return f.router.Register("OPTIONS", p, h)
}

// Register for TRACE APIs
func (f *fresh) Trace(p string, h HandlerFunc) error {
	return f.router.Register("TRACE", p, h)
}
