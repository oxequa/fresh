package fresh

import (
	"net/http"
	"io"
)

type Handler struct {

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "test")
}
