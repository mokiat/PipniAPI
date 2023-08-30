package registrymodel

var _ Resource = (*Workflow)(nil)

type Workflow struct {
	id   string
	name string
}

func (w *Workflow) ID() string {
	return w.id
}

func (w *Workflow) Name() string {
	return w.name
}
