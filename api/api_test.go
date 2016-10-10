package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
		token:    "JWT",
		endpoint: ts.URL,
	}

	v := url.Values{}
	v.Set("key", ExpectedTestAPIpost)

	r, err := a.post(CREATE, v)
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
