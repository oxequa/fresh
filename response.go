package fresh

import (
	"encoding/json"
	"net/http"
)

// Response structure
type (
	Response interface {
		Write(int, interface{})
	}

	response struct {
		w http.ResponseWriter
	}
)

// Response constructor
func NewResponse(w http.ResponseWriter) Response {
	w.Header().Set("Content-Type", "application/json")
	return &response{w: w}
}

// Response writer
func (r *response) Write(c int, i interface{}) {
	d, err := json.Marshal(i)
	if err != nil {
		return
	}
	r.w.WriteHeader(c)
	r.w.Write(d)
}
