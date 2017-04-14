package fresh

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
)

type (
	Fresh interface {
		Run()
		Get(string, func())
	}
	fresh struct {
		host    string
		port    string
		service *service // must be an array
	}
)

func New(h string, p string) Fresh {
	return &fresh{
		host: h,
		port: p,
		service: &service{
			server:  new(http.Server),
			handler: new(Handler),
			router:  new(Router),
		},
	}
	// config server array by reading JSON files fresh.json
}

func (f *fresh) Run() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.host+":"+f.port)
	if err != nil {
		os.Exit(1)
	}
	go func() {
		log.Println("Server started on " + f.host + ":" + f.port)
		f.service.server.Handler = f.service.handler
		f.service.server.Serve(listener)
	}()
	<-shutdown
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	f.service.server.Shutdown(ctx)
}

func (f *fresh) Get(p string, h func()) {
	// instantiate new route
	// append new route to router
}