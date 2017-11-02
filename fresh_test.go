package fresh

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"io"
	"encoding/json"
)

type r struct {
	path string
	args []string
}

var routes = []r{
	{
		path: "/first/",
		args: []string{},
	},
	{
		path: "/first/:/",
		args: []string{"1"},
	},
	{
		path: "/first/:/second",
		args: []string{"1"},
	},
	{
		path: "/first/:/second/:/",
		args: []string{"1","2"},
	},
	{
		path: "/first/:/second/:/third",
		args: []string{"1","2"},
	},
	{
		path: "/first/:/second/:/third/:/",
		args: []string{"1","2","3"},
	},
}

func ctrl(args []string) HandlerFunc{
	return func(context Context) error {
		data := make(map[string]string)
		for _, val := range args {
			data[val] = val
		}
		return context.Response().JSON(http.StatusOK, data)
	}
}

func setup() (fresh){
	f := fresh{
		config: new(config),
		server: new(http.Server),
		router: &router{&route{}, &context{}},
	}
	return f
}

func requests(method string, f *fresh){
	for _, elm := range routes {
		switch method {
		case "GET":
			f.GET(elm.path,ctrl(elm.args))
		case "POST":
			f.POST(elm.path,ctrl(elm.args))
		case "PUT":
			f.PUT(elm.path,ctrl(elm.args))
		case "TRACE":
			f.TRACE(elm.path,ctrl(elm.args))
		case "PATCH":
			f.PATCH(elm.path,ctrl(elm.args))
		case "DELETE":
			f.DELETE(elm.path,ctrl(elm.args))
		case "OPTIONS":
			f.OPTIONS(elm.path,ctrl(elm.args))
		}
	}
}

func records(method string, body io.Reader, f fresh, t *testing.T){
	for _, elm := range routes {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(method, elm.path, body)
		if err != nil {
			t.Fatal("Creating", method, elm, "request failed!")
		}
		f.router.ServeHTTP(rec,req)
		if rec.Code != http.StatusOK {
			t.Fatal("Server error: Returned ", rec.Code, " instead of ", http.StatusOK)
		}
		data := make(map[string]string)
		err = json.NewDecoder(rec.Body).Decode(&data)
		if err != nil {
			t.Error(err)
		}
		for _, value := range elm.args {
			if data[value] != value {
				t.Error("Response mismatch")
			}
		}
	}
}

func TestFresh_Run(t *testing.T) {
	f := New()
	go func() {
		err := f.Run()
		if err != nil {
			t.Error(err)
		}
		err = f.Shutdown()
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestFresh_GET(t *testing.T) {
	f := setup()
	requests("GET", &f)
	records("GET",nil, f, t)
}