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

type config struct {
	port  int    `json:"port,omitempty"`
	host  string `json:"host,omitempty"`
	ssl   bool   `json:"ssl,omitempty"`
	limit limit  `json:"limit,omitempty"` // body limit
	gzip  gzip   `json:"gzip,omitempty"`  // gzip config
	cors  cors   `json:"cors,omitempty"`  // cors options
}

type limit struct {
	status bool   `json:"status,omitempty"`
	size   string `json:"size,omitempty"`
}

type cors struct {
	status      bool   `json:"status,omitempty"`
	origins     string `json:"origins,omitempty"`
	headers     string `json:"headers,omitempty"`
	methods     string `json:"methods,omitempty"`
	credentials string `json:"credentials,omitempty"`
	expose      string `json:"expose,omitempty"`
	age         string `json:"age,omitempty"`
}

type gzip struct {
	status bool `json:"status,omitempty"`
	level  int  `json:"level,omitempty"`
}

func (c *config) read() error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &c)
}

func (c *config) write(path string) error {
	content, err := json.Marshal(c)
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
