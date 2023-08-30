package registrymodel

type Container interface {
	ID() string
	Name() string
	Resources() []Resource
	FindResource(id string) Resource
}

var _ Container = (*standardContainer)(nil)

const (
	RootContainerID = "f8802f97-4e1f-4fac-94c7-9fd02e2e3681"
)

type standardContainer struct {
	id   string
	name string

	children  []Container
	resources []Resource
}

func (c *standardContainer) ID() string {
	return c.id
}

func (c *standardContainer) Name() string {
	return c.name
}

func (c *standardContainer) Resources() []Resource {
	return c.resources
}

func (r *standardContainer) FindResource(id string) Resource {
	for _, resource := range r.resources {
		if resource.ID() == id {
			return resource
		}
	}
	for _, collection := range r.children {
		resource := collection.FindResource(id)
		if resource != nil {
			return resource
		}
	}
	return nil
}
