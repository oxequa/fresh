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
		SetTSL() Config
		SetPort(int) Config
		SetGzip(Gzip) Config
		SetDebug(bool) Config
		SetHost(string) Config
		SetLogger(bool) Config
		SetCertTSL(string, string) Config
	}

	config struct {
		*fresh
		Logs    bool        `json:"logs,omitempty"`    // srv lead
		Port    int         `json:"port,omitempty"`    // srv port
		Host    string      `json:"host,omitempty"`    // srv host
		Debug   bool        `json:"debug,omitempty"`   // debug status
		Logger  bool        `json:"logger,omitempty"`  // logger status
		TSL     *tls.Config `json:"tsl,omitempty"`     // tsl status
		Request *request    `json:"request,omitempty"` // request config
		Gzip    *Gzip       `json:"gzip,omitempty"`    // gzip config
		CORS    *cors       `json:"cors,omitempty"`    // cors options
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

func (c *config) SetTSL() Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.Host),
		Cache:      autocert.DirCache(".certs"), //folder for storing certificates
	}
	c.server.TLSConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}
	return c
}

func (c *config) SetGzip(g Gzip) Config {
	c.Gzip = &g
	return c
}

func (c *config) SetPort(port int) Config {
	// check if available
	c.Port = port
	return c
}

func (c *config) SetHost(host string) Config {
	// check if available
	c.Host = host
	return c
}

func (c *config) SetDebug(status bool) Config {
	c.Debug = status
	return c
}

func (c *config) SetLogger(status bool) Config {
	c.Logger = status
	return c
}

func (c *config) SetCertTSL(certFile, keyFile string) Config {
	return c
}
