package fresh

import (
	"testing"
)

func TestRun(t *testing.T) {
	f := New()
	f.Run()
	go func() {
		f.Run()
	}()
}
