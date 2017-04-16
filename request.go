package fresh

import (
	"net/http"
)


// Request structure
type (
	Request interface {

	}

	request struct {
		r *http.Request
	}
)

// Request constructor
func NewRequest(r *http.Request) Request{
	return &request{r: r}
}