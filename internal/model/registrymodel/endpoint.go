package registrymodel

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

func (e *Endpoint) Container() Container {
	return e.container
}
