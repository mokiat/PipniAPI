package view

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var EndpointManagement = co.Define(&endpointManagementComponent{})

type endpointManagementComponent struct {
	co.BaseComponent
}

func (c *endpointManagementComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Horizontal(layout.HorizontalSettings{
				ContentAlignment: layout.VerticalAlignmentCenter,
				ContentSpacing:   2,
			}),
		})

		co.WithChild("add", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/add.png"),
			})
		}))

		co.WithChild("edit", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/edit.png"),
			})
		}))

		co.WithChild("delete", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/delete.png"),
			})
		}))

		co.WithChild("separator", co.New(std.Spacing, func() {
			co.WithData(std.SpacingData{
				Size: ui.NewSize(10, 0),
			})
		}))

		co.WithChild("move-up", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/move-up.png"),
			})
		}))

		co.WithChild("move-down", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/move-down.png"),
			})
		}))
	})
}
