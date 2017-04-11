package fresh

import (
	"io"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "test")
	// find the right route that match current request
	// ger route and payload parameters
	// call route handler
}
