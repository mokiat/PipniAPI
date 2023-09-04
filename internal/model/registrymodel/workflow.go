package registrymodel

import (
	"fmt"

	"github.com/google/uuid"
)

var _ Resource = (*Workflow)(nil)

type Workflow struct {
	id        string
	name      string
	container Container
}

func (w *Workflow) ID() string {
	return w.id
}

func (w *Workflow) Name() string {
	return w.name
}

func (w *Workflow) SetName(name string) {
	w.name = name
}

func (w *Workflow) Kind() ResourceKind {
	return ResourceKindWorkflow
}

func (w *Workflow) Container() Container {
	return w.container
}

func (w *Workflow) Clone() Resource {
	return &Workflow{
		id:        uuid.Must(uuid.NewRandom()).String(),
		name:      fmt.Sprintf("%s Copy", w.name),
		container: w.container,
	}
}
