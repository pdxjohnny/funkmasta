package funkmasta

import (
	"io"
)

type Service struct {
	Name     string
	Endpoint string
	Payload  string
	EnvSetup string
}

type Request struct {
	Env    []string
	Body   io.Reader
	Output io.Writer
}
