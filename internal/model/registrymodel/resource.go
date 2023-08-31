package registrymodel

type Resource interface {
	ID() string
	Name() string
	Container() Container
}
