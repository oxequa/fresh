package fresh

import (
	"bytes"
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
	Config struct {
		*fresh   `yaml:"-"`
		request  *request          `yaml:"-"`                 // request config
		handlers []HandlerFunc     `yaml:"-"`                 // handlers array
		Host     string            `yaml:"host,omitempty"`    // server host
		Port     int               `yaml:"port,omitempty"`    // server port
		Logs     bool              `yaml:"logs,omitempty"`    // server logs
		Debug    bool              `yaml:"debug,omitempty"`   // debug status
		TSL      *TSL              `yaml:"tsl,omitempty"`     // tsl options
		Gzip     *Gzip             `yaml:"gzip,omitempty"`    // gzip Config
		CORS     *CORS             `yaml:"cors,omitempty"`    // cors options
		Limit    *Limit            `yaml:"limit,omitempty"`   // limit options
		Default  []string          `yaml:"default,omitempty"` // default static files (index.html or main.html and so on)
		Static   map[string]string `yaml:"static,omitempty"`  // serve static files
	}

	Limit struct {
		Body   string `yaml:"body,omitempty"`
		Header string `yaml:"header,omitempty"`
	}

	Gzip struct {
		writer         io.Writer
		responseWriter http.ResponseWriter
		Level          int      `yaml:"level,omitempty"`
		MinSize        int      `yaml:"size,omitempty"`
		Types          []string `yaml:"types,omitempty"`
		Filters        []Filter `yaml:"-,omitempty"`
	}

	TSL struct {
		Force bool   `yaml:"force,omitempty"`
		Crt   string `yaml:"crt,omitempty"`
		Key   string `yaml:"key,omitempty"`
	}

	CORS struct {
		Origins     []string `yaml:"origins,omitempty"`
		Methods     []string `yaml:"methods,omitempty"`
		Expose      []string `yaml:"expose,omitempty"`
		MaxAge      int      `yaml:"maxage,omitempty"`
		Credentials bool     `yaml:"credentials,omitempty"`
		Filters     []Filter `yaml:"-,omitempty"`
	}

	Security struct {
		XSS            string `yaml:"xss,omitempty"`
		HSTS           int    `yaml:"hsts,omitempty"`
		XDNS           bool   `yaml:"x-dns,omitempty"`
		CSFR           string `yaml:"csfr,omitempty"`
		XFrame         string `yaml:"x-frame,omitempty"`
		XContentType   string `yaml:"x-content-type,omitempty"`
		ReferrerPolicy string `yaml:"referrer-policy,omitempty"`
	}

	Filter func(Context) bool
)

// Init server config
func (c *Config) Init() *Config {
	// add handlers
	c.handlers = append(c.handlers,
		// gzip
		func(context Context) error {
			if c.Gzip != nil {
				reply := context.Response().get()
				// check buffer length
				if len(reply.response) >= c.Gzip.MinSize {
					r := context.Request().Get()
					w := context.Response().Get()
					if strings.Contains(r.Header.Get(AcceptEncoding), MIMEGzip) {
						ct := r.Header.Get(ContentType)
						if len(ct) == 0 || contain(ct, c.Gzip.Types) {
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
							if c.Gzip.Level >= gzip.NoCompression && c.Gzip.Level <= gzip.BestCompression {
								var err error
								gz, err = gzip.NewWriterLevel(w, c.Gzip.Level)
								if err != nil {
									context.Writer(w)
									return err
								}
							} else {
								gz = gzip.NewWriter(w)
							}
							context.Writer(&Gzip{writer: gz, responseWriter: w})
						}
					}
				}
			}
			return nil
		},
		// cors
		func(context Context) error {
			if c.CORS != nil {
				w := context.Response().Get()
				// Allow origins
				if len(c.CORS.Origins) > 0 {
					for _, h := range c.CORS.Origins {
						if h == "*" {
							w.Header().Set(AccessControlAllowOrigin, h)
							break
						} else if h == context.Request().Get().Header.Get("Origin") {
							w.Header().Set(AccessControlAllowOrigin, h)
						}
					}
				}
				// Allowed Methods
				if len(c.CORS.Methods) > 0 {
					w.Header().Set(AccessControlAllowMethods, strings.Join(c.CORS.Methods[:], ","))
				}
				// Allow credentials
				if c.CORS.Credentials {
					w.Header().Set(AccessControlAllowCredentials, "true")
				}
				// Expose headers
				if len(c.CORS.Expose) > 0 {
					w.Header().Set(AccessControlExposes, strings.Join(c.CORS.Expose[:], ","))
				}
				// Max age
				if c.CORS.MaxAge > 0 {
					w.Header().Set(AccessControlMaxAge, strconv.Itoa(c.CORS.MaxAge))
				}
			}
			return nil
		},
		// tsl
		func(context Context) error {
			return nil
		},
		// limit
		func(context Context) error {
			if c.Limit != nil {
				req := context.Request().Get()
				// content length
				if req.ContentLength > size(c.Limit.Body) {
					return nil
				}
				// read body
				buf := new(bytes.Buffer)
				l, err := buf.ReadFrom(req.Body)
				if err != nil || l > size(c.Limit.Body) {
					return nil
				}
			}
			return nil
		},
		// static
		func(context Context) error {
			return nil
		},
	)
	return c
}

// TSL with auto cert file
func (c *Config) tsl() *Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.Host),
	}
	c.Server.TLSConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}
	return c
}

// Read server Config from a file
func (c *Config) read(path string) error {
	content, err := ioutil.ReadFile(filepath.Join(path, file))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, &c)
}

// Write save server Config in a file
func (c *Config) write(path string) error {
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

// WriteHeader set a gzip header
func (g *Gzip) WriteHeader(i int) {
	g.responseWriter.WriteHeader(i)
}

// Header return gzip header
func (g *Gzip) Header() http.Header {
	return g.responseWriter.Header()
}

// Write gzip
func (g *Gzip) Write(b []byte) (int, error) {
	// check buffer
	return g.writer.Write(b)
}
