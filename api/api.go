package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type API struct {
	token    string
	endpoint string
}

func (a *API) post(resource string, data url.Values) (*http.Response, error) {
	// data := url.Values{}
	// data.Set("name", "foo")
	// data.Add("surname", "bar")

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
