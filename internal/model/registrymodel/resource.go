package registrymodel

type Resource interface {
	ID() string
	Name() string
	SetName(name string)
	Kind() ResourceKind
	Container() Container
}

type ResourceKind string

const (
	ResourceKindEndpoint ResourceKind = "endpoint"
	ResourceKindWorkflow ResourceKind = "workflow"
)
