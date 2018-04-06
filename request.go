package fresh

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

// Request structure
type (
	Request interface {
		setRouteParam(map[string]string)

		IsWS() bool
		IsTSL() bool
		URL() *url.URL
		Method() string
		JSON(interface{}) error
		Form() url.Values
		Get() *http.Request
		WS() *websocket.Conn
		Body() io.ReadCloser
		QueryString() string
		SetWS(*websocket.Conn)
		RouteParam(string) string
		FormValue(string) string
		QueryParam(string) string

		DTO
	}

	request struct {
		*context
		r  *http.Request
		ws *websocket.Conn
		p  map[string]string
		*dto
	}
)

// IsTSL check for a web socket request
func (req *request) IsWS() bool {
	h := req.r.Header.Get(Upgrade)
	return h == "websocket" || h == "Websocket"
}

// IsTSL check for a tsl request
func (req *request) IsTSL() bool {
	if req.r.TLS != nil {
		return true
	}
	return false
}

// Set URL parameters
func (req *request) URL() *url.URL {
	return req.r.URL
}

// Method current request
func (req *request) Method() string {
	return req.r.Method
}

// Get the form from a application/x-www-form-urlencoded request
func (req *request) Form() url.Values {
	return req.r.Form
}

// Get the body mapped to an interface from a application/json request
func (req *request) JSON(i interface{}) error {
	err := json.NewDecoder(req.r.Body).Decode(i)
	if err != nil {
		return err
	}
	// TODO: handle errors
	return nil
}

// Request return current http request
func (req *request) Get() *http.Request {
	return req.r
}

// IsTSL check for a web socket request
func (req *request) WS() *websocket.Conn {
	return req.ws
}

// Get the body from a application/json request
func (req *request) Body() io.ReadCloser {
	return req.r.Body
}

// Get the query string
func (req *request) QueryString() string {
	return req.r.URL.RawQuery
}

// SetWS used by the current request
func (req *request) SetWS(ws *websocket.Conn) {
	req.ws = ws
}

// Get a URL parameter
func (req *request) RouteParam(k string) string {
	return req.p[k]
}

// Get the form value by a given key from a application/x-www-form-urlencoded request
func (req *request) FormValue(k string) string {
	return req.r.FormValue(k)
}

// Get a query string parameter
func (req *request) QueryParam(k string) string {
	return req.r.URL.Query().Get(k)
}

// Set URL parameters
func (req *request) setRouteParam(m map[string]string) {
	req.p = m
}
