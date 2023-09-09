package registry

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var _ Resource = (*Endpoint)(nil)

type Endpoint struct {
	id        string
	name      string
	container Container

	method  string
	uri     string
	headers http.Header
	body    string
}

func (e *Endpoint) ID() string {
	return e.id
}

func (e *Endpoint) Name() string {
	return e.name
}

func (e *Endpoint) SetName(name string) {
	e.name = name
}

func (e *Endpoint) Kind() ResourceKind {
	return ResourceKindEndpoint
}

func (e *Endpoint) Container() Container {
	return e.container
}

func (e *Endpoint) Method() string {
	return e.method
}

func (e *Endpoint) SetMethod(method string) {
	e.method = method
}

func (e *Endpoint) URI() string {
	return e.uri
}

func (e *Endpoint) SetURI(uri string) {
	e.uri = uri
}

func (e *Endpoint) Headers() http.Header {
	return e.headers.Clone()
}

func (e *Endpoint) SetHeaders(headers http.Header) {
	e.headers = headers
}

func (e *Endpoint) Body() string {
	return e.body
}

func (e *Endpoint) SetBody(body string) {
	e.body = body
}

func (e *Endpoint) Clone() Resource {
	return &Endpoint{
		id:        uuid.Must(uuid.NewRandom()).String(),
		name:      fmt.Sprintf("%s Copy", e.name),
		container: e.container,

		method:  e.method,
		uri:     e.uri,
		headers: e.headers.Clone(),
		body:    e.body,
	}
}
