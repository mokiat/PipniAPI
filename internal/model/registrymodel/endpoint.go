package registrymodel

var _ Resource = (*Endpoint)(nil)

type Endpoint struct {
	id   string
	name string
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
