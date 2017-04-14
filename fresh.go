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
		Run() error
		Get(string, func()) error
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

func (f *fresh) Run() error {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.host+":"+f.port)
	if err != nil {
		return err
	}
	go func() {
		log.Println("Server started on " + f.host + ":" + f.port)
		f.service.server.Handler = f.service.handler
		f.service.server.Serve(listener)
	}()
	<-shutdown
	log.Println("Server shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	f.service.server.Shutdown(ctx)
	return nil
}

func (f *fresh) Get(p string, h func()) error{
	// instantiate new route
	// append new route to router
}