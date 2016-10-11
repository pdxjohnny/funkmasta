package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pdxjohnny/getfunky/getfunky"
)

const (
	ExpectedTestAPIpost = "value"
)

func TestAPIpost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := readerToString(r.Body)
		r.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		v, err := url.ParseQuery(b)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprint(w, v.Get("key"))
	}))
	defer ts.Close()

	a := &API{
		endpoint: ts.URL,
	}

	v := url.Values{}
	v.Set("key", ExpectedTestAPIpost)

	r, err := a.post("/test", v)
	if err != nil {
		t.Fatal(err)
	}

	b, err := readerToString(r.Body)
	r.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if ExpectedTestAPIpost != b {
		t.Fatalf("Expected: %v, got: %v", ExpectedTestAPIpost, b)
	}
}

func TestAPICreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	}))
	defer ts.Close()

	a := &API{
		endpoint: ts.URL,
	}

	err := a.Create(&getfunky.Service{
		Name:     "name",
		Endpoint: "name.service",
		Payload:  `def test():\n  print("Hello World")`,
		Env:      "virtualenv .venv\n. .venv/bin/activate",
	})
	if err != nil {
		t.Fatal(err)
	}
}
