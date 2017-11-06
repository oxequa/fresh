package fresh

import (
	"crypto/tls"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	file = "fresh.json"
	perm = 0770
)

type (
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

	config struct {
		*fresh
		logs          bool        `json:"logs,omitempty"`           // srv lead
		port          int         `json:"port,omitempty"`           // srv port
		host          string      `json:"host,omitempty"`           // srv host
		debug         bool        `json:"debug,omitempty"`          // debug status
		logger        bool        `json:"logger,omitempty"`         // logger status
		tsl           *tls.Config `json:"tsl,omitempty"`            // tsl status
		request       *request    `json:"request,omitempty"`        // request config
		gzip          *Gzip       `json:"gzip,omitempty"`           // gzip config
		cors          *cors       `json:"cors,omitempty"`           // cors options
		staticDefault []string    `json:"static_default,omitempty"` // default static files served
	}

	limits struct {
		BodyLimit   string `json:"body_limit,omitempty"`
		HeaderLimit string `json:"header_limit,omitempty"`
	}

	cors struct {
		Status      bool   `json:"status,omitempty"`
		Origins     string `json:"origins,omitempty"`
		Headers     string `json:"headers,omitempty"`
		Methods     string `json:"methods,omitempty"`
		Credentials string `json:"credentials,omitempty"`
		Expose      string `json:"expose,omitempty"`
		Age         string `json:"age,omitempty"`
	}

	Gzip struct {
		Status  bool     `json:"status,omitempty"`
		Level   int      `json:"level,omitempty"`
		MinSize int      `json:"size,omitempty"`
		Types   []string `json:"types,omitempty"`
	}
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
	c.gzip = &g
	return c
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
