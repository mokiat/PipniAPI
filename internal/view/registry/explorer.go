package registry

import (
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Explorer = mvc.EventListener(co.Define(&explorerComponent{}))

type ExplorerData struct {
	RegistryModel *registry.Model
}

type explorerComponent struct {
	co.BaseComponent

	mdlRegistry *registry.Model
}

func (c *explorerComponent) OnUpsert() {
	data := co.GetData[ExplorerData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *explorerComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.SymmetricSpacing(0, 1),
			BorderColor: opt.V(std.OutlineColor),
			Layout:      layout.Fill(),
			Padding:     ui.SymmetricSpacing(0, 5),
		})

		co.WithChild("list", co.New(std.List, func() {

			for _, resource := range c.mdlRegistry.Root().Resources() {
				resource := resource
				co.WithChild(resource.ID(), co.New(Item, func() {
					co.WithData(ItemData{
						Selected: c.mdlRegistry.SelectedID() == resource.ID(),
						Icon:     c.resourceImage(resource),
						Text:     resource.Name(),
					})
					co.WithCallbackData(ItemCallbackData{
						OnClick: func() {
							c.onResourceSelected(resource)
						},
					})
				}))
			}
		}))
	})
}

func (c *explorerComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case registry.RegistrySelectionChangedEvent:
		c.Invalidate()
	case registry.RegistryStructureChangedEvent:
		c.Invalidate()
	}
}

func (c *explorerComponent) onResourceSelected(resource registry.Resource) {
	c.mdlRegistry.SetSelectedID(resource.ID())
}

func (c *explorerComponent) resourceImage(resource registry.Resource) *ui.Image {
	switch resource.(type) {
	case *registry.Endpoint:
		return co.OpenImage(c.Scope(), "images/ping.png")
	case *registry.Workflow:
		return co.OpenImage(c.Scope(), "images/workflow.png")
	default:
		return nil
	}
}
