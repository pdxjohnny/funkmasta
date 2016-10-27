package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pdxjohnny/getfunky/getfunky"
)

// API is an http.Client associated with an endpoint, e.g. getfunky.example.com
// the API can access services of that endpoint or create, update, and delete
// services on that encpoint
type API struct {
	http.Client
	token    string
	endpoint string
}

// Get wraps http.Client.Get so that url is the path from the endpoint
func (a *API) Get(path string) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path

	return a.Client.Get(u.String())
}

// Head wraps http.Client.Get so that url is the path from the endpoint
func (a *API) Head(path string) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path

	return a.Client.Head(u.String())
}

// Post wraps http.Client.Get so that url is the path from the endpoint
func (a *API) Post(path string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path

	return a.Client.Post(u.String(), bodyType, body)
}

// PostForm wraps http.Client.Get so that url is the path from the endpoint
func (a *API) PostForm(path string, data url.Values) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path

	return a.Client.PostForm(u.String(), data)
}

// PostService is used to interact with (call) created services
func (a *API) PostService(service string, data url.Values, body io.Reader) (*http.Response, error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = service
	u.RawQuery = data.Encode()

	res, err := a.Client.Post(u.String(), "application/octet-stream", body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create tells the endpoint to create a new service with the data from
// getfunky.Service
func (a *API) Create(s *getfunky.Service) error {
	v := url.Values{}
	v.Set("name", s.Name)
	v.Set("endpoint", s.Endpoint)
	v.Set("payload", s.Payload)
	v.Set("env", s.EnvSetup)

	_, err := a.PostForm(CREATE, v)
	if err != nil {
		return err
	}

	return nil
}

func ParseCreate(r io.Reader) (*getfunky.Service, error) {
	qBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	q := string(qBytes)

	v, err := url.ParseQuery(q)
	if err != nil {
		return nil, err
	}

	s := &getfunky.Service{
		Name:     v.Get("name"),
		Endpoint: v.Get("endpoint"),
		Payload:  v.Get("payload"),
		EnvSetup: v.Get("env"),
	}

	return s, nil
}
