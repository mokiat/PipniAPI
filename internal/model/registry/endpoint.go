package registry

import (
	"fmt"

	"github.com/google/uuid"
)

var _ Resource = (*Endpoint)(nil)

type Endpoint struct {
	id        string
	name      string
	container Container
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

func (e *Endpoint) Clone() Resource {
	return &Endpoint{
		id:        uuid.Must(uuid.NewRandom()).String(),
		name:      fmt.Sprintf("%s Copy", e.name),
		container: e.container,
		// TODO: Copy more stuff here
	}
}
