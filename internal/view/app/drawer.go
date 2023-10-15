package app

import (
	"github.com/mokiat/PipniAPI/internal/model/registry"
	contextview "github.com/mokiat/PipniAPI/internal/view/context"
	registryview "github.com/mokiat/PipniAPI/internal/view/registry"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Drawer = co.Define(&drawerComponent{})

type DrawerData struct {
	RegistryModel *registry.Model
}

type drawerComponent struct {
	co.BaseComponent

	mdlRegistry *registry.Model
}

func (c *drawerComponent) OnUpsert() {
	data := co.GetData[DrawerData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *drawerComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderColor: opt.V(std.OutlineColor),
			BorderSize: ui.Spacing{
				Right: 1,
			},
			Padding: ui.UniformSpacing(5),
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Top:    5,
					Bottom: 5,
				},
			}),
		})

		co.WithChild("environment-selection", co.New(contextview.Selector, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(contextview.SelectorData{
				RegistryModel: c.mdlRegistry,
			})
		}))

		co.WithChild("endpoint-selection", co.New(registryview.Explorer, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(registryview.ExplorerData{
				RegistryModel: c.mdlRegistry,
			})
		}))

		co.WithChild("endpoint-management", co.New(registryview.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentBottom,
			})
			co.WithData(registryview.ToolbarData{
				RegistryModel: c.mdlRegistry,
			})
		}))
	})
}
