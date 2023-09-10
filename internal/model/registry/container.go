package registry

import "golang.org/x/exp/slices"

type Container interface {
	ID() string
	Name() string
	Resources() []Resource
	AppendResource(resource Resource)
	RemoveResource(resource Resource)
	FindResource(id string) Resource
	ResourcePosition(resource Resource) int
	MoveResourceUp(resource Resource)
	MoveResourceDown(resource Resource)
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

func (c *standardContainer) AppendResource(resource Resource) {
	c.resources = append(c.resources, resource)
}

func (c *standardContainer) RemoveResource(resource Resource) {
	c.resources = slices.DeleteFunc(c.resources, func(candidate Resource) bool {
		return candidate == resource
	})
}

func (c *standardContainer) FindResource(id string) Resource {
	for _, resource := range c.resources {
		if resource.ID() == id {
			return resource
		}
	}
	for _, collection := range c.children {
		resource := collection.FindResource(id)
		if resource != nil {
			return resource
		}
	}
	return nil
}

func (c *standardContainer) ResourcePosition(resource Resource) int {
	for i, candidate := range c.resources {
		if candidate == resource {
			return i
		}
	}
	return -1
}

func (c *standardContainer) MoveResourceUp(resource Resource) {
	position := c.ResourcePosition(resource)
	if position <= 0 {
		return
	}
	c.swapResourcesAtPositions(position-1, position)
}

func (c *standardContainer) MoveResourceDown(resource Resource) {
	position := c.ResourcePosition(resource)
	if position >= len(c.resources)-1 {
		return
	}
	c.swapResourcesAtPositions(position, position+1)
}

func (c *standardContainer) swapResourcesAtPositions(first, second int) {
	c.resources[first], c.resources[second] = c.resources[second], c.resources[first]
}
