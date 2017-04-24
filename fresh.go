package fresh

import (
	"context"
	"log"
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
		Get(string, func(Request, Response)) error
		Post(string, func(Request, Response)) error
		Put(string, func(Request, Response)) error
		Patch(string, func(Request, Response)) error
		Delete(string, func(Request, Response)) error
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
		config: &config{
			Host: "localhost",
			Port: 3000,
		},
	}
	wd, _ := os.Getwd()
	if fresh.config.read(wd) != nil {
		// create a config with default params
		fresh.config.write(wd)
		return &fresh
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
