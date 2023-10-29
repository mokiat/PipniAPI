package registry

type Resource interface {
	ID() string
	Name() string
	SetName(name string)
	Kind() ResourceKind
	Container() Container
	Clone() Resource
}

type ResourceKind string

const (
	ResourceKindEndpoint ResourceKind = "endpoint"
	ResourceKindWorkflow ResourceKind = "workflow"
	ResourceKindContext  ResourceKind = "context"
)
