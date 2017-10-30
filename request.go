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
		setURLParam(map[string] string)

		IsWS() bool
		IsTSL() bool
		Method() string
		Map(interface{})
		Form() url.Values
		Get() *http.Request
		WS() *websocket.Conn
		Body() io.ReadCloser
		QueryString() string
		SetWS(*websocket.Conn)
		URLParam(string) string
		FormValue(string) string
		QueryParam(string) string
	}

	request struct {
		*context
		r *http.Request
		ws *websocket.Conn
		p map[string] string
	}
)

// Set URL parameters
func (req *request) setURLParam(m map[string] string) {
	req.p = m
}

// IsTSL check for a web socket request
func (req *request) IsWS() bool{
	h := req.r.Header.Get(Upgrade)
	return h == "websocket" || h == "Websocket"
}

// IsTSL check for a tsl request
func (req *request) IsTSL() bool{
	if req.r.TLS != nil{
		return true
	}
	return false
}

// Method current request
func (req *request) Method() string{
	return req.r.Method
}

// Get the form from a application/x-www-form-urlencoded request
func (req *request) Form() url.Values {
	return req.r.Form
}

// Get the body mapped to an interface from a application/json request
func (req *request) Map(i interface{}) {
	err := json.NewDecoder(req.r.Body).Decode(i)
	if err != nil {
		return
	}
	// TODO: handle errors
}

// Request return current http request
func (req *request) Get() *http.Request {
	return req.r
}

// IsTSL check for a web socket request
func (req *request) WS() *websocket.Conn{
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
func (req *request) SetWS(ws *websocket.Conn){
	req.ws = ws
}

// Get a URL parameter
func (req *request) URLParam(k string) string {
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
