package fresh

import (
	"net/http"
	"net"
	"os"
	"fmt"
)

type (
	Fresh interface {
		New(string, string)
		Run()
	}
	fresh struct {
		host string
		port string
		server *http.Server
		handler *Handler
		router *Router
	}
)

func New(h string, p string) *fresh{
	return &fresh{host: h, port: p, server:new(http.Server), handler: new(Handler)}
}

func (f *fresh) Run(){
	listener, error := net.Listen("tcp", f.host + ":" + f.port)
	if error != nil{
		os.Exit(1)
	}
	fmt.Println("Server started on " + f.host + ":" + f.port)
	f.server.Handler = f.handler
	f.server.Serve(listener)
}




