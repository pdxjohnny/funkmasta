package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pdxjohnny/getfunky/getfunky"
)

const (
	expectedTestAPIpost = "value"

	expectedTestAPICreateName             = "name"
	expectedTestAPICreateEndpoint         = "service.name"
	expectedTestAPICreatePayloadPlaintext = "def test():\n  print(\"Hello World\")"
	expectedTestAPICreatePayloadBinary    = "\xbd\xb2\x3d\x00\xFF\xbc\x20\xe2\x8c\x98"
	expectedTestAPICreateEnvSetup         = "virtualenv .venv\n. .venv/bin/activate"
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
	v.Set("key", expectedTestAPIpost)

	r, err := a.post("/test", v)
	if err != nil {
		t.Fatal(err)
	}

	b, err := readerToString(r.Body)
	r.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if expectedTestAPIpost != b {
		t.Fatalf("Expected: %v, got: %v", expectedTestAPIpost, b)
	}
}

func TestAPICreatePlaintext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := ParseCreate(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if s.Name != expectedTestAPICreateName {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateName, s.Name)
		} else if s.Endpoint != expectedTestAPICreateEndpoint {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateEndpoint, s.Endpoint)
		} else if s.Payload != expectedTestAPICreatePayloadPlaintext {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreatePayloadPlaintext, s.Payload)
		} else if s.EnvSetup != expectedTestAPICreateEnvSetup {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateEnvSetup, s.EnvSetup)
		}
	}))
	defer ts.Close()

	a := &API{
		endpoint: ts.URL,
	}

	err := a.Create(&getfunky.Service{
		Name:     expectedTestAPICreateName,
		Endpoint: expectedTestAPICreateEndpoint,
		Payload:  expectedTestAPICreatePayloadPlaintext,
		EnvSetup: expectedTestAPICreateEnvSetup,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPICreateBinary(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := ParseCreate(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if s.Name != expectedTestAPICreateName {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateName, s.Name)
		} else if s.Endpoint != expectedTestAPICreateEndpoint {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateEndpoint, s.Endpoint)
		} else if s.Payload != expectedTestAPICreatePayloadBinary {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreatePayloadBinary, s.Payload)
		} else if s.EnvSetup != expectedTestAPICreateEnvSetup {
			t.Fatalf("Expected: %v, got: %v", expectedTestAPICreateEnvSetup, s.EnvSetup)
		}
	}))
	defer ts.Close()

	a := &API{
		endpoint: ts.URL,
	}

	err := a.Create(&getfunky.Service{
		Name:     expectedTestAPICreateName,
		Endpoint: expectedTestAPICreateEndpoint,
		Payload:  expectedTestAPICreatePayloadBinary,
		EnvSetup: expectedTestAPICreateEnvSetup,
	})
	if err != nil {
		t.Fatal(err)
	}
}
