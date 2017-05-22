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
	"fmt"
)

// Main Fresh structure
type (
	Fresh interface {
		Run() error
		Get(string, func(Request, Response)) error
		Post(string, func(Request, Response)) error
		Put(string, func(Request, Response)) error
		Trace(string, func(Request, Response)) error
		Patch(string, func(Request, Response)) error
		Delete(string, func(Request, Response)) error
		Options(string, func(Request, Response)) error
	}
	fresh struct {
		config *config
		router *Router
		server *http.Server
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

// Set a list of http header fields
func (f *fresh) Header(){

}

// Load all servers configurations and start them
func (f *fresh) Run() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.config.Host+":"+strconv.Itoa(f.config.Port))
	if err != nil {
		fmt.Println(err)
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
func (f *fresh) Get(p string, h func(Request, Response)) error {
	return f.router.Register("GET", p, h)
}

// Register for POST APIs
func (f *fresh) Post(p string, h func(Request, Response)) error {
	return f.router.Register("POST", p, h)
}

// Register for PUT APIs
func (f *fresh) Put(p string, h func(Request, Response)) error {
	return f.router.Register("PUT", p, h)
}

// Register for PATCH APIs
func (f *fresh) Patch(p string, h func(Request, Response)) error {
	return f.router.Register("PATCH", p, h)
}

// Register for DELETE APIs
func (f *fresh) Delete(p string, h func(Request, Response)) error {
	return f.router.Register("DELETE", p, h)
}

// Register for OPTIONS APIs
func (f *fresh) Options(p string, h func(Request, Response)) error {
	return f.router.Register("OPTIONS", p, h)
}

// Register for TRACE APIs
func (f *fresh) Trace(p string, h func(Request, Response)) error {
	return f.router.Register("TRACE", p, h)
}
