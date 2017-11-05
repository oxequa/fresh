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
		writeErr(error)
		Code(int) error
		Type(content string)
		Raw(int, string) error
		Error(int, error) error
		HTML(int, string) error
		File(int, string) error
		Get() http.ResponseWriter
		Download(int, string) error
		XML(int, interface{}) error
		Text(int, interface{}) error
		JSON(int, interface{}) error
		JSONP(int, string, interface{}) error
		XMLFormat(int, interface{}, string) error
		JSONFormat(int, interface{}, string) error
		JSONPFormat(int, string, interface{}, string) error
	}

	Reply struct {
		code     int
		response []byte
	}

	response struct {
		*context
		w     http.ResponseWriter
		r     *http.Request
		reply Reply
	}
)

// Return writer
func (r *response) write() {
	r.w.WriteHeader(r.reply.code)
	r.w.Write(r.reply.response)
}

func (r *response) writeErr(err error) {
	if r.reply.code == 0 {
		http.Error(r.w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(r.w, err.Error(), r.reply.code)
	}
}

// Check content type
func (r *response) check(content string) {
	if r.w.Header().Get(ContentType) != content {
		r.w.Header().Set(ContentType, content)
	}
}

// Set response values
func (r *response) set(code int, response []byte) {
	r.reply = Reply{code, response}
}

// Http code response
func (r *response) Code(c int) error {
	r.set(c, nil)
	return nil
}

// Type set content type
func (r *response) Type(content string) {
	r.w.Header().Set(ContentType, content)
}

// Http code response
func (r *response) Get() http.ResponseWriter {
	return r.w
}

// Raw response
func (r *response) Raw(c int, i string) error {
	r.set(c, []byte(i))
	return nil
}

// HTTP error
func (r *response) Error(c int, err error) error {
	r.set(c, nil)
	return err
}

// HTML response
func (r *response) HTML(c int, i string) error {
	r.check(MIMETextHTML)
	r.set(c, []byte(i))
	return nil
}

// File response, may be used to display a file
func (r *response) File(c int, path string) error {
	// check if exist
	if _, err := os.Stat(path); err != nil {
		r.set(http.StatusNotFound, nil)
		return err
	}
	// check if is a dir or not
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		r.set(http.StatusNotFound, nil)
		return err
	}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	http.ServeContent(r.w, r.r, fi.Name(), fi.ModTime(), f)
	return nil
}

// XML response
func (r *response) XML(c int, i interface{}) error {
	r.check(MIMEAppXML)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, nil)
		return err
	}
	r.set(c, d)
	return nil
}

// JSON response
func (r *response) JSON(c int, i interface{}) error {
	r.check(MIMEAppJSON)
	d, err := json.Marshal(i)
	if err != nil {
		r.set(c, nil)
		return err
	}
	r.set(c, d)
	return nil
}

// Text response
func (r *response) Text(c int, i interface{}) error {
	r.check(MIMEText)
	d, err := xml.Marshal(i)
	if err != nil {
		r.set(c, nil)
		return err
	}
	r.set(c, d)
	return nil
}

// Download response, force the browser to download a file
func (r *response) Download(c int, path string) error {
	// check if exist
	if _, err := os.Stat(path); err != nil {
		r.set(http.StatusNotFound, nil)
		return err
	}
	// check if is a dir or not
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		r.set(http.StatusNotFound, nil)
		return err
	}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	r.w.Header().Set("Content-Disposition", "attachment; filename="+fi.Name())
	http.ServeContent(r.w, r.r, fi.Name(), fi.ModTime(), f)
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
		r.set(c, nil)
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)))
	return nil
}

// XML pretty response
func (r *response) XMLFormat(c int, i interface{}, indent string) error {
	r.check(MIMEAppXML)
	d, err := xml.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, nil)
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
		r.set(c, nil)
		return err
	}
	r.set(c, d)
	return nil
}

// JSON pretty response
func (r *response) JSONPFormat(c int, callback string, i interface{}, indent string) error {
	r.check(MIMEAppJS)
	d, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		r.set(c, nil)
		return err
	}
	r.set(c, []byte(fmt.Sprintf("%s(%s)", callback, d)))
	return nil
}
