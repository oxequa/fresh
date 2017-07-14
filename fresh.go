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

// MIME types
const (
	MIMEAppJSON    = "application/json" + ";" + UTF8
	MIMEAppJS      = "application/javascript" + ";" + UTF8
	MIMEAppXML     = "application/xml" + ";" + UTF8
	MIMEUrlencoded = "application/x-www-form-urlencoded"
	MIMEMultipart  = "multipart/form-data"
	MIMETextHTML   = "text/html" + ";" + UTF8
	MIMETextXML    = "text/xml" + ";" + UTF8
	MIMEText       = "text/plain" + ";" + UTF8
)

// Encoding Chartset
const (
	UTF8     = "charset=UTF-8"
	ISO88591 = "chartset=ISO-8859-1"
)

// Header types
const (
	Location                      = "Location"
	ContentType                   = "Content-Type"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// Main Fresh structure
type (
	Fresh interface {
		Run() error
		Get(string, HandlerFunc) Handler
		Post(string, HandlerFunc) Handler
		Put(string, HandlerFunc) Handler
		Trace(string, HandlerFunc) Handler
		Patch(string, HandlerFunc) Handler
		Delete(string, HandlerFunc) Handler
		Options(string, HandlerFunc) Handler

		After(...HandlerFunc) Fresh
		Before(...HandlerFunc) Fresh
		Group(string) Fresh
		Rest(string, ...HandlerFunc) Resource
	}

	fresh struct {
		group  *route
		config *config
		router *router
		server *http.Server
	}

	Context interface {
		Request() Request
		Response() Response
	}

	context struct {
		request  Request
		response Response
	}

	HandlerFunc func(Context) error
)

// Initialize main Fresh structure
func New() Fresh {
	fresh := fresh{
		config: new(config),
		server: new(http.Server),
		router: new(router),
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

// Set context request and response
func (c *context) new(r *http.Request, w http.ResponseWriter) {
	c.response = &response{w: w, r: r}
	c.request = &request{r: r}
}

// Return context request
func (c *context) Request() Request {
	return c.request
}

// Return context response
func (c *context) Response() Response {
	return c.response
}

// Register a group
func (f fresh) Group(path string) Fresh {
	f.group = &route{
		path: strings.Split(path, "/"),
	}
	return &f
}

// After middleware
func (f *fresh) After(middleware ...HandlerFunc) Fresh {
	f.group.after = append(f.group.after, middleware...)
	return f
}

// Before middleware
func (f *fresh) Before(middleware ...HandlerFunc) Fresh {
	f.group.before = append(f.group.before, middleware...)
	return f
}

// Register for GET APIs
func (f *fresh) Get(path string, handler HandlerFunc) Handler {
	return f.router.register("GET", path, f.group, handler)
}

// Register for PUT APIs
func (f *fresh) Put(path string, handler HandlerFunc) Handler {
	return f.router.register("PUT", path, f.group, handler)
}

// Register for POST APIs
func (f *fresh) Post(path string, handler HandlerFunc) Handler {
	return f.router.register("POST", path, f.group, handler)
}

// Register for TRACE APIs
func (f *fresh) Trace(path string, handler HandlerFunc) Handler {
	return f.router.register("TRACE", path, f.group, handler)
}

// Register for PATCH APIs
func (f *fresh) Patch(path string, handler HandlerFunc) Handler {
	return f.router.register("PATCH", path, f.group, handler)
}

// Register for DELETE APIs
func (f *fresh) Delete(path string, handler HandlerFunc) Handler {
	return f.router.register("DELETE", path, f.group, handler)
}

// Register for OPTIONS APIs
func (f *fresh) Options(path string, handler HandlerFunc) Handler {
	return f.router.register("OPTIONS", path, f.group, handler)
}

// Register a resource (get, post, put, delete)
func (f *fresh) Rest(path string, handlers ...HandlerFunc) Resource {
	res := resource{
		methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}
	path = strings.Trim(path, "/")
	name := "{" + path + "}"
	if strings.LastIndex(path, "/") != -1 {
		name = string("{" + path[strings.LastIndex(path, "/")+1:] + "}")
	}
	for _, method := range res.methods {
		switch method {
		case "GET":
			res.rest = append(res.rest, f.router.register(method, path, f.group, handlers[0]))
		case "POST":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, f.group, handlers[1]))
		case "PUT", "PATCH":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, f.group, handlers[2]))
		case "DELETE":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, f.group, handlers[3]))
		}
	}
	return &res
}
