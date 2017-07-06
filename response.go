package fresh

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

type (
	Response interface {
		write()
		Code(int) error
		Raw(int, string) error
		Text(int, interface{}) error
		HTML(int, string) error
		XML(int, interface{}) error
		XMLFormat(int, interface{}, string) error
		JSON(int, interface{}) error
		JSONFormat(int, interface{}, string) error
		JSONP(int, string, interface{}) error
		JSONPFormat(int, string, interface{}, string) error
	}

	response struct {
		w        http.ResponseWriter
		err      error
		code     int
		response []byte
	}
)

// Response constructor
func newResponse(writer http.ResponseWriter) Response {
	return &response{w: writer}
}

// Check content type
func (r *response) check(content string) {
	if r.w.Header().Get(ContentType) == "" {
		r.w.Header().Set(ContentType, content)
	}
}

// Return writer
func (r *response) write() {
	if r.code != 0 && r.response != nil {
		r.w.WriteHeader(r.code)
		r.w.Write(r.response)
	}
}

// Set response values
func (r *response) set(code int, response []byte) {
	r.response = response
	r.code = code
}

// Custom Header
func (r *response) Header(key, value string) Response {
	r.w.Header().Add(key, value)
	return r
}

// Text response
func (r *response) Text(c int, i interface{}) error {
	r.check(MIMEText)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, d)
	return nil
}

// Http code response
func (r *response) Code(c int) error {
	r.set(c, []byte{})
	return nil
}

// Raw response
func (r *response) Raw(c int, i string) error {
	r.set(c, []byte(i))
	return nil
}

// JSON response
func (r *response) JSON(c int, i interface{}) error {
	r.check(MIMEAppJSON)
	d, err := json.Marshal(i)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, d)
	return nil
}

// JSON pretty response
func (r *response) JSONFormat(c int, i interface{}, indent string) error {
	r.check(MIMEAppJSON)
	d, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, d)
	return nil
}

// JSON response
func (r *response) JSONP(c int, callback string, i interface{}) error {
	r.check(MIMEAppJS)
	d, err := json.Marshal(i)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)))
	return nil
}

// JSON pretty response
func (r *response) JSONPFormat(c int, callback string, i interface{}, indent string) error {
	r.check(MIMEAppJS)
	d, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)))
	return nil
}

// XML response
func (r *response) XML(c int, i interface{}) error {
	r.check(MIMEAppXML)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, d)
	return nil
}

// XML pretty response
func (r *response) XMLFormat(c int, i interface{}, indent string) error {
	r.check(MIMEAppXML)
	d, err := xml.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, []byte(err.Error()))
		return err
	}
	r.set(c, d)
	return nil
}

// HTML response
func (r *response) HTML(c int, i string) error {
	r.check(MIMETextHTML)
	r.set(c, []byte(i))
	return nil
}

// File response

// Download response

// Redirect
func (r *response) Redirect(c int, link string) error {
	if c < http.StatusMultipleChoices || c > http.StatusTemporaryRedirect {
		return nil
	}
	r.w.Header().Set(Location, link)
	r.w.WriteHeader(c)
	return nil
}
