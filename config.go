package fresh

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v2"
)

const (
	file = "fresh.yaml"
	perm = 0770
)

type (
	config struct {
		*fresh
		port          int           // server port
		host          string        // server host
		logs          bool          // server lead
		debug         bool          // debug status
		logger        bool          // logger status
		tsl           *tls.Config   // tsl status
		request       *request      // request config
		gzip          *Gzip         // gzip config
		cors          *CORS         // cors options
		handlers      []HandlerFunc // handlers array
		staticDefault []string      // default static files served
	}

	Limit struct {
		BodyLimit   string `yaml:"body_limit,omitempty"`
		HeaderLimit string `yaml:"header_limit,omitempty"`
	}

	Gzip struct {
		writer         io.Writer
		responseWriter http.ResponseWriter
		Level          int      `yaml:"level,omitempty"`
		MinSize        int      `yaml:"size,omitempty"`
		Types          []string `yaml:"types,omitempty"`
		Filter         Filter
	}

	CORS struct {
		Origins     []string `yaml:"origins,omitempty"`
		Methods     []string `yaml:"methods,omitempty"`
		Expose      []string `yaml:"expose,omitempty"`
		MaxAge      int      `yaml:"maxage,omitempty"`
		Credentials bool     `yaml:"credentials,omitempty"`
		Filter      Filter
	}

	Config interface {
		TSL() Config
		Port(int) Config
		CORS(CORS) Config
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
	return yaml.Unmarshal(content, &c)
}

func (c *config) write(path string) error {
	content, err := yaml.Marshal(c)
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
		w := context.Response().Get()
		// Allow origins
		if len(c.cors.Origins) > 0 {
			for _, h := range c.cors.Origins {
				if h == "*" {
					w.Header().Set(AccessControlAllowOrigin, h)
					break
				} else if h == context.Request().Get().Header.Get("Origin") {
					w.Header().Set(AccessControlAllowOrigin, h)
				}
			}
		}
		// Allowed Methods
		if len(c.cors.Methods) > 0 {
			w.Header().Set(AccessControlAllowMethods, strings.Join(c.cors.Methods[:], ","))
		}
		// Allow credentials
		if c.cors.Credentials {
			w.Header().Set(AccessControlAllowCredentials, "true")
		}
		// Expose headers
		if len(c.cors.Expose) > 0 {
			w.Header().Set(AccessControlExposes, strings.Join(c.cors.Expose[:], ","))
		}
		// Max age
		if c.cors.MaxAge > 0 {
			w.Header().Set(AccessControlMaxAge, strconv.Itoa(c.cors.MaxAge))
		}
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
