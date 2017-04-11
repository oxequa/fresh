package fresh

import (
	"io/ioutil"
	"os"
	"testing"
	"encoding/json"
	"path/filepath"
)

func temp() (d string, err error){
	d, err = ioutil.TempDir(os.TempDir(),dir)
	return
}

func TestRead(t *testing.T) {
	config := &config{}
	content, err := json.Marshal(config)
	if err != nil{
		t.Error(err)
	}
	path, err := temp()
	if err != nil{
		t.Error(err)
	}
	ioutil.WriteFile(filepath.Join(path,file), content, perm)
	if err != nil {
		t.Error(err)
	}
	if err = config.read(path); err != nil{
		t.Error(err)
	}
	os.Remove(path)
}

func TestWrite(t *testing.T){
	config := &config{}
	path, err := temp()
	if err != nil{
		t.Error(err)
	}
	if err = config.write(path); err != nil{
		t.Error(err)
	}
}
