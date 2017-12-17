package fresh

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	file = "fresh.json"
	perm = 0770
)

type (
	config struct {
		*fresh
		logs          bool          // srv lead
		port          int           // srv port
		host          string        // srv host
		debug         bool          // debug status
		logger        bool          // logger status
		tsl           *tls.Config   // tsl status
		request       *request      // request config
		gzip          *Gzip         // gzip config
		cors          *CORS         // cors options
		handlers      []HandlerFunc // handlers array
		staticDefault []string      // default static files served
	}

	limits struct {
		BodyLimit   string `json:"body_limit,omitempty"`
		HeaderLimit string `json:"header_limit,omitempty"`
	}

	Gzip struct {
		writer         io.Writer
		responseWriter http.ResponseWriter
		Level          int      `json:"level,omitempty"`
		MinSize        int      `json:"size,omitempty"`
		Types          []string `json:"types,omitempty"`
		Filter         Filter
	}

	CORS struct {
		Origins     string `json:"origins,omitempty"`
		Headers     string `json:"headers,omitempty"`
		Methods     string `json:"methods,omitempty"`
		Credentials string `json:"credentials,omitempty"`
		Expose      string `json:"expose,omitempty"`
		MaxAge      string `json:"maxage,omitempty"`
		Filter      Filter
	}

	Config interface {
		TSL() Config
		Port(int) Config
		Gzip(Gzip) Config
		Debug(bool) Config
		Host(string) Config
		Logger(bool) Config
		CertTSL(string, string) Config
		StaticDefault([]string) Config
	}

	Filter func(Context) bool
)

func (c *config) read(path string) error {
	content, err := ioutil.ReadFile(filepath.Join(path, file))
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &c)
}

func (c *config) write(path string) error {
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, perm); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filepath.Join(path, file), content, perm)
}

func (c *config) append(handler HandlerFunc) Config {
	c.handlers = append(c.handlers, handler)
	return c
}

func (c *config) contains(s string, arr []string) bool {
	s = strings.ToLower(s)
	for _, val := range arr {
		if val == s {
			return true
		}
	}
	return false
}

func (c *config) TSL() Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.host),
		Cache:      autocert.DirCache(".certs"), //folder for storing certificates
	}
	c.server.TLSConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}
	return c
}

func (c *config) Gzip(g Gzip) Config {
	// set only value not nil or zero or not set
	c.gzip = &g
	// gzip handler
	handler := func(context Context) (err error) {
		reply := context.Response().get()
		// check buffer length
		if len(reply.response) >= c.gzip.MinSize {
			r := context.Request().Get()
			w := context.Response().Get()
			if strings.Contains(r.Header.Get(AcceptEncoding), MIMEGzip) {
				ct := r.Header.Get(ContentType)
				if len(ct) == 0 || c.contains(ct, c.gzip.Types) {
					if len(ct) == 0 {
						// detect content type by reading response
						w.Header().Set(ContentType, http.DetectContentType(reply.response))
					}
					// set header
					w.Header().Set(ContentEncoding, MIMEGzip)
					// del length if exist
					w.Header().Del(ContentLength)
					// new writer
					gz := &gzip.Writer{}
					defer gz.Close()
					if c.gzip.Level >= gzip.NoCompression && c.gzip.Level <= gzip.BestCompression {
						gz, err = gzip.NewWriterLevel(w, c.gzip.Level)
						if err != nil {
							context.Writer(w)
							return err
						}
					} else {
						gz = gzip.NewWriter(w)
					}
					context.Writer(Gzip{writer: gz, responseWriter: w})
				}
			}
		}
		return nil
	}
	return c.append(handler)
}

func (c *config) CORS(s CORS) Config {
	c.cors = &s
	// cors handler
	handler := func(context Context) error {
		return nil
	}
	return c.append(handler)
}

func (c *config) Port(port int) Config {
	// check if available
	c.port = port
	return c
}

func (c *config) Host(host string) Config {
	// check if available
	c.host = host
	return c
}

func (c *config) Debug(status bool) Config {
	c.debug = status
	return c
}

func (c *config) Logger(status bool) Config {
	c.logger = status
	return c
}

func (c *config) CertTSL(certFile, keyFile string) Config {
	return c
}

func (c *config) StaticDefault(staticDefault []string) Config {
	c.staticDefault = staticDefault
	return c
}

func (g Gzip) WriteHeader(i int) {
	g.responseWriter.WriteHeader(i)
}

func (g Gzip) Header() http.Header {
	return g.responseWriter.Header()
}

func (g Gzip) Write(b []byte) (int, error) {
	// check buffer
	return g.writer.Write(b)
}
