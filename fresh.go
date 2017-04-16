package fresh

import (
	"fmt"
	"net"
	"net/http"
	"os"
)


// Main Fresh structure
type (
	Fresh interface {
		Run()
		Get(string, func(Request, Response))
	}
	fresh struct {
		host    string
		port    string
		service *service // must be an array
	}
)


// Initialize main Fresh structure
func New(h string, p string) Fresh {
	return &fresh{
		host: h,
		port: p,
		service: &service{
			server:  new(http.Server),
			router:  new(Router),
		},
	}
	// config server array by reading JSON files fresh.json
}


// Load all servers configurations and start them
func (f *fresh) Run() {
	listener, err := net.Listen("tcp", f.host + ":" + f.port)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Server started on " + f.host + ":" + f.port)
	f.service.server.Handler = f.service.router
	f.service.server.Serve(listener)
}


// Register for GET APIs
func (f *fresh) Get(p string, h func(Request, Response)) {
	r := &Route{
		method:	"GET",
		path: p,
		handler: h}
	f.service.router.routes = append(f.service.router.routes, r)
}


