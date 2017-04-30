package fresh

import (
	"testing"
)

func TestRun(t *testing.T) {
	f := New()
	go func() {
		f.Run()
	}()
}
