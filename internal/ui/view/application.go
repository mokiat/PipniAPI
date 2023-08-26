package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Application = co.Define(&applicationComponent{})

type applicationComponent struct {
	co.BaseComponent
}

func (c *applicationComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Frame(),
		})

		co.WithChild("toolbar", co.New(Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(ToolbarData{})
		}))

		co.WithChild("content", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ElementData{
				Layout: layout.Frame(),
			})

			co.WithChild("drawer", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(300),
				})
				co.WithData(std.ContainerData{
					BackgroundColor: opt.V(ui.Red()),
					BorderColor:     opt.V(std.OutlineColor),
					BorderSize: ui.Spacing{
						Right: 1,
					},
				})
			}))

			co.WithChild("workspace", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ContainerData{
					BackgroundColor: opt.V(ui.Green()),
				})
			}))

		}))
	})
}
