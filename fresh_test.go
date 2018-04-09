package fresh

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testRoute struct {
	path      string
	urlParams []string
	body      dataResponse
	code      int
}

type dataResponse struct {
	Field  int      `json:"field"`
	Fields []string `json:"fields"`
}

var req = []testRoute{
	{
		path:      "/first/",
		urlParams: []string{},
		body:      dataResponse{Field: 1, Fields: []string{"white", "black"}},
		code:      http.StatusOK,
	},
	{
		path:      "/first/:/",
		urlParams: []string{"1"},
		body:      dataResponse{Field: 1, Fields: []string{"brown", "charcoal"}},
		code:      http.StatusCreated,
	},
	{
		path:      "/first/:/second",
		urlParams: []string{"1"},
		body:      dataResponse{Field: 1, Fields: []string{"yellow"}},
		code:      http.StatusOK,
	},
	{
		path:      "/first/:/second/:/",
		urlParams: []string{"1", "2"},
		body:      dataResponse{Field: 1, Fields: []string{"blue"}},
		code:      http.StatusForbidden,
	},
	{
		path:      "/first/:/second/:/third",
		urlParams: []string{"1", "2"},
		body:      dataResponse{Field: 1, Fields: []string{"red"}},
		code:      http.StatusOK,
	},
	{
		path:      "/first/:/second/:/third/:/",
		urlParams: []string{"1", "2", "3"},
		body:      dataResponse{Field: 1, Fields: []string{"blue", "violet"}},
		code:      http.StatusOK,
	},
}

func setup() fresh {
	f := fresh{
		Config: &Config{},
		Server: new(http.Server),
	}
	f.Config.fresh = &f
	f.Config.Init()
	f.router = &router{&f, &route{}, make(map[string]string)}
	return f
}

func ctrl(r testRoute) HandlerFunc {
	return func(context Context) error {
		return context.Response().JSON(r.code, r.body)
	}
}

func requests(method string, f *fresh) {
	for _, elm := range req {
		// set url params
		for _, v := range elm.urlParams {
			elm.path = strings.Replace(elm.path, ":/", ":"+v+"/", 1)
		}
		switch method {
		case "GET":
			f.GET(elm.path, ctrl(elm))
		case "POST":
			f.POST(elm.path, ctrl(elm))
		case "PUT":
			f.PUT(elm.path, ctrl(elm))
		case "TRACE":
			f.TRACE(elm.path, ctrl(elm))
		case "PATCH":
			f.PATCH(elm.path, ctrl(elm))
		case "DELETE":
			f.DELETE(elm.path, ctrl(elm))
		case "OPTIONS":
			f.OPTIONS(elm.path, ctrl(elm))
		}
	}
}

func records(method string, body io.Reader, f fresh, t *testing.T) {
	for _, elm := range req {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(method, elm.path, body)
		if err != nil {
			t.Fatal("Creating", method, elm, "request failed!")
		}
		f.router.ServeHTTP(rec, req)
		// status code
		if rec.Code != elm.code {
			t.Fatal("Server error: Returned ", rec.Code, " instead of ", elm.code)
		}
		// body
		expected, _ := json.Marshal(elm.body)
		result, _ := ioutil.ReadAll(rec.Body)
		if string(result) != string(expected) {
			t.Fatal("Expected", string(expected), "instead", string(result))
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

	// Test multiple calls
	//TODO test
	// rec := httptest.NewRecorder()
	// f.GET("multiple", func(c Context) error {
	// 	time.Sleep(5 * time.Second)
	// 	return c.Response().JSON(2, nil)
	// })
	// var wg sync.WaitGroup
	// wg.Add(10)
	// for i := 1; i <= 10; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		req, err := http.NewRequest("GET", "multiple", nil)
	// 		if err != nil {
	// 			t.Fatal("Creating GET multiple request failed!")
	// 		}
	// 		fresh.router.ServeHTTP(rec, req)
	// 		// status code
	// 		if rec.Code != 200 {
	// 			t.Fatal("Server error: Returned ", rec.Code, " instead of ", 200)
	// 		}
	// 	}()
	// }
	// wg.Wait()
}

func TestFresh_GET(t *testing.T) {
	f := setup()
	requests("GET", &f)
	records("GET", nil, f, t)
}

func TestFresh_PUT(t *testing.T) {
	f := setup()
	requests("PUT", &f)
	records("PUT", nil, f, t)
}

func TestFresh_POST(t *testing.T) {
	f := setup()
	requests("POST", &f)
	records("POST", nil, f, t)
}

func TestFresh_TRACE(t *testing.T) {
	f := setup()
	requests("TRACE", &f)
	records("TRACE", nil, f, t)
}

func TestFresh_PATCH(t *testing.T) {
	f := setup()
	requests("PATCH", &f)
	records("PATCH", nil, f, t)
}

func TestFresh_DELETE(t *testing.T) {
	f := setup()
	requests("DELETE", &f)
	records("DELETE", nil, f, t)
}

func TestFresh_OPTIONS(t *testing.T) {
	f := setup()
	requests("OPTIONS", &f)
	records("OPTIONS", nil, f, t)
}
