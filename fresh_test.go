package fresh

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const PORT = 8080

type Data struct {
	Data string `json:"data"`
}

var REQUESTS = map[string][]string{
	"/first/":                    {""},
	"/first/:/":                  {"1"},
	"/first/:/second/":           {"1"},
	"/first/:/second/:/":         {"1", "2"},
	"/first/:/second/:/third/":   {"1", "2"},
	"/first/:/second/:/third/:/": {"1", "2", "3"},
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
	f := New()
	f.Config().SetPort(PORT)
	for key, values := range REQUESTS {
		data := make(map[string]string)
		for _, val := range values {
			data[val] = val
		}
		f.GET(makeRequestURL(key, values, true), func(c Context) error {
			return c.Response().JSON(http.StatusOK, data)
		})
	}
	go func() {
		f.Run()
	}()
	for key, values := range REQUESTS {
		resp, err := http.Get(makeRequestURL(key, values, false))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Error(resp.StatusCode)
		}
		data := make(map[string]string)
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			t.Error(err)
		}
		for _, value := range values {
			if data[value] != value {
				t.Error("Response mismatch")
			}
		}
	}
	f.Shutdown()
}

func TestFresh_POST(t *testing.T) {
	f := New()
	f.Config().SetPort(PORT)
	for key, values := range REQUESTS {
		data := make(map[string]string)
		for _, val := range values {
			data[val] = val
		}
		f.POST(makeRequestURL(key, values, true), func(c Context) error {
			data["body"] = c.Request().FormValue("body")
			return c.Response().JSON(http.StatusOK, data)
		})
	}
	go func() {
		f.Run()
	}()
	for key, values := range REQUESTS {
		body := map[string]string{"body": "body"}
		jsonBody, _ := json.Marshal(body)
		resp, err := http.Post(makeRequestURL(key, values, false), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Error(resp.StatusCode)
		}
		data := make(map[string]string)
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			t.Error(err)
		}
		for _, value := range values {
			if data[value] != value {
				t.Error("Response mismatch")
			}
		}
	}
	f.Shutdown()
}

func makeRequestURL(url string, values []string, parameter bool) string {
	sep := ""
	host := "http://localhost:" + strconv.Itoa(PORT)
	if parameter {
		sep = ":"
		host = ""
	}
	for _, value := range values {
		url = strings.Replace(url, ":/", sep+value+"/", 1)
	}
	return host + url
}
