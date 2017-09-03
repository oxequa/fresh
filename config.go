package fresh

import (
	"encoding/json"
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
		SetPort(int) Config
		SetHost(string) Config
		SetTSL(string, string) Config
	}

	config struct {
		Port    int    `json:"port,omitempty"` // srv port
		Host    string `json:"host,omitempty"` // srv host
		TSL     *tsl   `json:"tsl,omitempty"`  // tsl status
		Request *request `json:"request,omitempty"` // request config
		Gzip    *gzip    `json:"gzip,omitempty"`    // gzip config
		CORS    *cors    `json:"cors,omitempty"`    // cors options
	}

	tsl struct {

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

	gzip struct {
		Status bool `json:"status,omitempty"`
		Level  int  `json:"level,omitempty"`
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

func (c *config) SetTSL(cert, key string) Config {
	return c
}

func (c *config) AutoTSL(cert, key string) Config {
	return c
}