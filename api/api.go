package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pdxjohnny/getfunky/getfunky"
)

type API struct {
	token    string
	endpoint string
}

func (a *API) post(resource string, data url.Values) (*http.Response, error) {
	u, err := url.ParseRequestURI(a.endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = resource

	encoded := data.Encode()
	req, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewBufferString(encoded),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", a.token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readerToString(r io.Reader) (string, error) {
	b := new(bytes.Buffer)
	l, err := b.ReadFrom(r)
	if err != nil {
		return "", err
	} else if l < 1 {
		return "", nil
	}

	return b.String(), nil
}

func (a *API) Create(s *getfunky.Service) error {
	v := url.Values{}
	v.Set("name", s.Name)
	v.Set("endpoint", s.Endpoint)
	v.Set("payload", s.Payload)
	v.Set("env", s.EnvSetup)

	_, err := a.post(CREATE, v)
	if err != nil {
		return err
	}

	return nil
}

func ParseCreate(r io.Reader) (*getfunky.Service, error) {
	q, err := readerToString(r)
	if err != nil {
		return nil, err
	}

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
