package srv

type Config struct {
	Ssl     bool     `json:"ssl,omitempty"`
	Request *request `json:"request,omitempty"` // request config
	Gzip    *gzip    `json:"gzip,omitempty"`    // gzip config
	Cors    *cors    `json:"cors,omitempty"`    // cors options
}

type request struct {
	BodyLimit   string `json:"body_limit,omitempty"`
	HeaderLimit string `json:",omitempty"`
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
