package context

type Environment struct {
	id   string
	name string
}

func (e *Environment) ID() string {
	return e.id
}

func (e *Environment) Name() string {
	return e.name
}
