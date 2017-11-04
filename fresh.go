package fresh

import (
	httpContext "context"
	"golang.org/x/net/websocket"
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

// Access
const (
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllows           = "Access-Control-Allow-s"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlExposes          = "Access-Control-Expose-s"
	AccessControlRequests         = "Access-Control-Request-s"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// Request
const (
	Accept              = "Accept"
	AcceptEncoding      = "Accept-Encoding"
	Allow               = "Allow"
	Authorization       = "Authorization"
	ContentDisposition  = "Content-Disposition"
	ContentEncoding     = "Content-Encoding"
	ContentLength       = "Content-Length"
	ContentType         = "Content-Type"
	Cookie              = "Cookie"
	SetCookie           = "Set-Cookie"
	IfModifiedSince     = "If-Modified-Since"
	LastModified        = "Last-Modified"
	Location            = "Location"
	Upgrade             = "Upgrade"
	Vary                = "Vary"
	WWWAuthenticate     = "WWW-Authenticate"
	XForwardedFor       = "X-Forwarded-For"
	XForwardedProto     = "X-Forwarded-Proto"
	XForwardedProtocol  = "X-Forwarded-Protocol"
	XForwardedSsl       = "X-Forwarded-Ssl"
	XUrlScheme          = "X-Url-Scheme"
	XHTTPMethodOverride = "X-HTTP-Method-Override"
	XRealIP             = "X-Real-IP"
	XRequestID          = "X-Request-ID"
	Server              = "Server"
	Origin              = "Origin"
)

// Encoding chartset
const (
	UTF8     = "charset=UTF-8"
	ISO88591 = "chartset=ISO-8859-1"
)

// Main Fresh structure
type (
	Fresh interface {
		Rest
		Run() error
		Shutdown() error
		Config() Config
		Group(string) Group
	}

	fresh struct {
		config *config
		router *router
		static map[string]string
		server *http.Server
	}

	Rest interface {
		STATIC(map[string]string)
		WS(string, HandlerFunc) Handler
		GET(string, HandlerFunc) Handler
		POST(string, HandlerFunc) Handler
		PUT(string, HandlerFunc) Handler
		TRACE(string, HandlerFunc) Handler
		PATCH(string, HandlerFunc) Handler
		DELETE(string, HandlerFunc) Handler
		OPTIONS(string, HandlerFunc) Handler
		CRUD(string, ...HandlerFunc) Resource
	}

	Context interface {
		Request() Request
		Response() Response
	}

	context struct {
		request    Request
		response   Response
		parameters map[string]string
	}

	HandlerFunc func(Context) error
)

// Initialize main Fresh structure
func New() Fresh {
	fresh := fresh{
		config: new(config),
		server: new(http.Server),
		router: &router{&route{}, make(map[string]string)},
	}
	wd, _ := os.Getwd()
	if fresh.config.read(wd) != nil {
		// random port
		rand.Seed(time.Now().Unix())
		fresh.config.Host = "127.0.0.1"
		fresh.config.Port = rand.Intn(9999-1111) + 1111
	}
	return &fresh
}

// Run HTTP server
func (f *fresh) Run() error {
	shutdown := make(chan os.Signal)
	port := strconv.Itoa(f.config.Port)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	listener, err := net.Listen("tcp", f.config.Host+":"+port)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go func() {
		log.Println("Server listen on", f.config.Host+":"+port)
		f.server.Handler = f.router
		f.router.printRoutes()
		// check for tsl before serve
		f.server.Serve(listener)
	}()
	<-shutdown
	log.Println("Server shutdown")
	ctx, cancel := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	cancel()
	return nil
}

// Config interface
func (f *fresh) Config() Config {
	return f.config
}

// Shutdown server
func (f *fresh) Shutdown() error {
	ctx, cancel := httpContext.WithTimeout(httpContext.Background(), 5*time.Second)
	f.server.Shutdown(ctx)
	cancel()
	return nil

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

// WS api registration
func (f *fresh) WS(path string, handler HandlerFunc) Handler {
	h := func(c Context) (err error) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			c.Request().SetWS(ws)
			err = handler(c)
		}).ServeHTTP(c.Response().Get(), c.Request().Get())
		return err
	}
	return f.router.register("GET", path, h)
}

// Register a resource (get, post, put, delete)
func (f *fresh) CRUD(path string, h ...HandlerFunc) Resource {
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
			res.rest = append(res.rest, f.router.register(method, path, h[0]))
		case "POST":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[1]))
		case "PUT", "PATCH":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[2]))
		case "DELETE":
			res.rest = append(res.rest, f.router.register(method, path+"/"+name, h[3]))
		}
	}
	return &res
}

// GET api registration
func (f *fresh) GET(path string, handler HandlerFunc) Handler {
	return f.router.register("GET", path, handler)
}

// PUT api registration
func (f *fresh) PUT(path string, handler HandlerFunc) Handler {
	return f.router.register("PUT", path, handler)
}

// POST api registration
func (f *fresh) POST(path string, handler HandlerFunc) Handler {
	return f.router.register("POST", path, handler)
}

// TRACE api registration
func (f *fresh) TRACE(path string, handler HandlerFunc) Handler {
	return f.router.register("TRACE", path, handler)
}

// PATCH api registration
func (f *fresh) PATCH(path string, handler HandlerFunc) Handler {
	return f.router.register("PATCH", path, handler)
}

// DELETE api registration
func (f *fresh) DELETE(path string, handler HandlerFunc) Handler {
	return f.router.register("DELETE", path, handler)
}

// OPTIONS api registration
func (f *fresh) OPTIONS(path string, handler HandlerFunc) Handler {
	return f.router.register("OPTIONS", path, handler)
}

// ASSETS serve a list of static files. Array of files or directories TODO write logic
func (f *fresh) STATIC(static map[string]string) {
	f.router.registerStatic(static)
}

// Init set context request and response
func (c *context) init(r *http.Request, w http.ResponseWriter) {
	c.response = &response{w: w, r: r}
	c.request = &request{r: r}
	c.request.setURLParam(c.parameters)
}

// Return context request
func (c *context) Request() Request {
	return c.request
}

// Return context response
func (c *context) Response() Response {
	return c.response
}
