package fresh

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const dir = "fresh"

func temp() (d string, err error) {
	d, err = ioutil.TempDir(os.TempDir(), dir)
	return
}

func TestRead(t *testing.T) {
	result := config{}
	expected := config{Port: 3000}

	content, err := json.Marshal(&expected)
	if err != nil {
		t.Error(err)
	}
	path, err := temp()
	if err != nil {
		t.Error(err)
	}
	ioutil.WriteFile(filepath.Join(path, file), content, perm)
	if err != nil {
		t.Error(err)
	}
	if err = result.read(path); err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Error("Expected:", expected, "instead", result)
	}
	os.Remove(path)
}

func TestWrite(t *testing.T) {
	config := config{}
	path, err := temp()
	if err != nil {
		t.Error(err)
	}
	if err = config.write(path); err != nil {
		t.Error(err)
	}
	os.Remove(path)
}
