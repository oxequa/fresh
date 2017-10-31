package fresh

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
)

type (
	Response interface {
		write()
		Code(int) error
		Get() http.ResponseWriter
		Raw(int, string) error
		HTML(int, string) error
		File(int, string) error
		Download(int, string) error
		XML(int, interface{}) error
		Text(int, interface{}) error
		JSON(int, interface{}) error
		JSONP(int, string, interface{}) error
		XMLFormat(int, interface{}, string) error
		JSONFormat(int, interface{}, string) error
		JSONPFormat(int, string, interface{}, string) error
	}

	response struct {
		*context
		w        http.ResponseWriter
		r        *http.Request
		err      error
		code     int
		response []byte
	}
)

// Return writer
func (r *response) write() {
	if r.code != 0 && r.response != nil {
		r.w.WriteHeader(r.code)
		r.w.Write(r.response)
	}
}

// Check content type
func (r *response) check(content string) {
	if r.w.Header().Get(ContentType) == "" {
		r.w.Header().Set(ContentType, content)
	}
}

// Set response values
func (r *response) set(code int, response []byte, err error) {
	r.response = response
	r.code = code
	r.err = err
}

// Http code response
func (r *response) Code(c int) error {
	r.set(c, nil, nil)
	return nil
}

// Http code response
func (r *response) Get() http.ResponseWriter {
	return r.w
}

// Raw response
func (r *response) Raw(c int, i string) error {
	r.set(c, []byte(i), nil)
	return nil
}

// HTML response
func (r *response) HTML(c int, i string) error {
	r.check(MIMETextHTML)
	r.set(c, []byte(i), nil)
	return nil
}

// File response, may be used to display a file
func (r *response) File(c int, path string) error {
	// check if exist
	if _, err := os.Stat(path); err != nil {
		r.set(http.StatusNotFound, nil, err)
		return nil
	}
	// check if is a dir or not
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		r.set(http.StatusNotFound, nil, err)
		return nil
	}
	f, _ := os.Open(path)
	http.ServeContent(r.w, r.r, fi.Name(), fi.ModTime(), f)
	f.Close()
	return nil
}

// XML response
func (r *response) XML(c int, i interface{}) error {
	r.check(MIMEAppXML)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, d, nil)
	return nil
}

// JSON response
func (r *response) JSON(c int, i interface{}) error {
	r.check(MIMEAppJSON)
	d, err := json.Marshal(i)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, d, nil)
	return nil
}

// Text response
func (r *response) Text(c int, i interface{}) error {
	r.check(MIMEText)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, d, nil)
	return nil
}

// Download response, force the browser to download a file
func (r *response) Download(c int, path string) error {
	// check if exist
	if _, err := os.Stat(path); err != nil {
		r.set(http.StatusNotFound, nil, err)
		return nil
	}
	// check if is a dir or not
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		r.set(http.StatusNotFound, nil, err)
		return nil
	}
	f, _ := os.Open(path)
	r.w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	http.ServeContent(r.w, r.r, fi.Name(), fi.ModTime(), f)
	f.Close()
	return nil
}

// Redirect
func (r *response) Redirect(c int, link string) error {
	if c < http.StatusMultipleChoices || c > http.StatusTemporaryRedirect {
		return nil
	}
	r.w.Header().Set(Location, link)
	r.w.WriteHeader(c)
	return nil
}

// Custom Header
func (r *response) Header(key, value string) Response {
	r.w.Header().Add(key, value)
	return r
}

// JSON response
func (r *response) JSONP(c int, callback string, i interface{}) error {
	r.check(MIMEAppJS)
	d, err := json.Marshal(i)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)), nil)
	return nil
}

// XML pretty response
func (r *response) XMLFormat(c int, i interface{}, indent string) error {
	r.check(MIMEAppXML)
	d, err := xml.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, d, nil)
	return nil
}

// JSON pretty response
func (r *response) JSONFormat(c int, i interface{}, indent string) error {
	r.check(MIMEAppJSON)
	d, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, d, nil)
	return nil
}

// JSON pretty response
func (r *response) JSONPFormat(c int, callback string, i interface{}, indent string) error {
	r.check(MIMEAppJS)
	d, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, nil, err)
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)), nil)
	return nil
}
