package registry

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/mokiat/gog"
)

var _ Resource = (*Endpoint)(nil)

type Endpoint struct {
	id        string
	name      string
	container Container

	method  string
	uri     string
	headers []gog.KV[string, string]
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

func (e *Endpoint) Headers() []gog.KV[string, string] {
	return slices.Clone(e.headers)
}

func (e *Endpoint) SetHeaders(headers []gog.KV[string, string]) {
	e.headers = slices.Clone(headers)
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
		headers: slices.Clone(e.headers),
		body:    e.body,
	}
}
