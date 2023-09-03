package registrymodel

type Resource interface {
	ID() string
	Name() string
	Container() Container
}

type ResourceKind string

const (
	ResourceKindEndpoint ResourceKind = "endpoint"
	ResourceKindWorkflow ResourceKind = "workflow"
)
