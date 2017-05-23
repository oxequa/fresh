package fresh

import (
	"encoding/json"
	"net/http"
)

type (
	Response interface {
		Write(int, interface{}) error
	}

	response struct {
		w http.ResponseWriter
	}
)

// Response constructor
func NewResponse(writer http.ResponseWriter) Response {
	writer.Header().Set("Content-Type", "application/json")
	return &response{w: writer}
}

// Response writer
func (r *response) Write(c int, i interface{}) error {
	d, err := json.Marshal(i)
	if err != nil {
		return err
	}
	r.w.WriteHeader(c)
	r.w.Write(d)
	return nil
}
