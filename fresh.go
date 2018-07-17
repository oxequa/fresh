package fresh

import (
	httpContext "context"
	"log"
	"net/http"
	"os"
	"time"

	"net"
	"os/signal"
	"strconv"
	"syscall"
)

// Main Fresh structure
type (
	fresh struct {
		config *Config
		router *router
		server *http.Server
	}

	context struct {
		request    request
		response   response
		parameters map[string]string
	}

	Fresh interface {
		Rest
		Stop() error
		Start() error
		Config() *Config
		Group(string) Group
	}

	Context interface {
		Request() Request
		Response() Response
		Writer(http.ResponseWriter)
	}

	HandlerFunc func(Context) error
)

// New Fresh instance
func New() Fresh {
	fresh := fresh{}
	// Fresh config
	fresh.config = &Config{}
	fresh.config.init(&fresh)
	// Server setting
	fresh.server = new(http.Server)
	// Fresh router
	fresh.router = &router{&fresh, &route{}, make(map[string]string)}

	wd, _ := os.Executable()
	if fresh.config.read(wd) != nil {
		fresh.config.Host = "127.0.0.1"
		fresh.config.Port = randPort(fresh.config.Host, 3000)
	}
	return &fresh
}

// Shutdown server
func (f *fresh) Stop() error {
	ctx, cancel := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	cancel()
	f.config.log("Server shutdown")
	return nil
}

// Start HTTP server
func (f *fresh) Start() error {
	shutdown := make(chan os.Signal)
	port := strconv.Itoa(f.config.Port)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	listener, err := net.Listen("tcp", f.config.Host+":"+port)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go func() {
		f.config.banner()
		f.config.log("Server listen on", f.config.Host+":"+port)
		f.server.Handler = f.router
		// default route
		if f.router.route.children == nil {
			f.GET("/", func(c Context) error {
				return c.Response().Raw(http.StatusOK, welcome)
			})
		}
		if f.config.Router != nil && f.config.Router.Print {
			PrintRouter(f.router)
		}
		// check for tsl before serve
		if f.config.TSL != nil {
			f.config.tsl()
		}
		f.server.Serve(listener)
	}()
	<-shutdown
	f.Stop()
	return nil
}

// Config return server settings
func (f *fresh) Config() *Config {
	return f.config
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
