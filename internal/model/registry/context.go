package registry

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/mokiat/gog"
)

var _ Resource = (*Context)(nil)

type Context struct {
	id         string
	name       string
	container  Container
	properties []gog.KV[string, string]
}

func (c *Context) ID() string {
	return c.id
}

func (c *Context) Name() string {
	return c.name
}

func (c *Context) SetName(name string) {
	c.name = name
}

func (c *Context) Kind() ResourceKind {
	return ResourceKindContext
}

func (c *Context) Container() Container {
	return c.container
}

func (c *Context) Properties() []gog.KV[string, string] {
	return slices.Clone(c.properties)
}

func (c *Context) SetProperties(properties []gog.KV[string, string]) {
	c.properties = slices.Clone(properties)
}

func (c *Context) Clone() Resource {
	return &Context{
		id:        uuid.Must(uuid.NewRandom()).String(),
		name:      fmt.Sprintf("%s Copy", c.name),
		container: c.container,

		properties: slices.Clone(c.properties),
	}
}
