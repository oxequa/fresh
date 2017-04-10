package fresh

import (
	"testing"
)

func TestRun(t *testing.T) {
	host := "localhost"
	port := "8080"
	f := New(host, port)
	f.Run()
}
