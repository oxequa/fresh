package fresh

import (
	"net/http"
	"net/url"
	"io"
	"encoding/json"
)


// Request structure
type (
	Request interface {
		//QueryString() string
		//QueryParam(string) string
		//Param(string) string
		Body() io.ReadCloser
		Map(interface{})
		Form() url.Values
		FormValue(string) string
	}

	request struct {
		r *http.Request
	}
)

// Request constructor
func NewRequest(r *http.Request) Request{
	return &request{r: r}
}


// Return the body from a application/json request
func(req * request) Body() io.ReadCloser{
	return req.r.Body
}

// Return the body mapped to an interface from a application/json request
func(req * request) Map(i interface{}){
	err := json.NewDecoder(req.r.Body).Decode(i)
	if err != nil{
		return
	}
	// TODO: handle errors
}

// Return the form from a application/x-www-form-urlencoded request
func(req * request) Form() url.Values{
	return req.r.Form
}

// Return the form value by a given key from a application/x-www-form-urlencoded request
func(req * request) FormValue(k string) string{
	return req.r.FormValue(k)
}