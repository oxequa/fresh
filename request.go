package fresh

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Request structure
type (
	Request interface {
		QueryString() string
		QueryParam(string) string
		URLParam(string) string
		Body() io.ReadCloser
		Map(interface{})
		Form() url.Values
		FormValue(string) string
	}

	request struct {
		httpRequest *http.Request
	}
)

// Request constructor
func newRequest(r *http.Request) Request {
	return &request{httpRequest: r}
}

// Get the query string
func (req *request) QueryString() string {
	return req.httpRequest.URL.RawQuery
}

// Get a query string parameter
func (req *request) QueryParam(k string) string {
	return req.httpRequest.URL.Query().Get(k)
}

// Get a URL parameter
func (req *request) URLParam(k string) string {
	return req.httpRequest.URL.Query().Get(k)
}

// Get the body from a application/json request
func (req *request) Body() io.ReadCloser {
	return req.httpRequest.Body
}

// Get the body mapped to an interface from a application/json request
func (req *request) Map(i interface{}) {
	err := json.NewDecoder(req.httpRequest.Body).Decode(i)
	if err != nil {
		return
	}
	// TODO: handle errors
}

// Get the form from a application/x-www-form-urlencoded request
func (req *request) Form() url.Values {
	return req.httpRequest.Form
}

// Get the form value by a given key from a application/x-www-form-urlencoded request
func (req *request) FormValue(k string) string {
	return req.httpRequest.FormValue(k)
}
