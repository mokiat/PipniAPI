package view

import (
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var EnvironmentSelection = co.Define(&environmentSelectionComponent{})

type environmentSelectionComponent struct {
	co.BaseComponent
}

func (c *environmentSelectionComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Right: 5,
				},
			}),
		})

		co.WithChild("dropdown", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.DropdownData{
				Items: []std.DropdownItem{
					{
						Key:   "staging",
						Label: "Staging",
					},
					{
						Key:   "production",
						Label: "Production",
					},
				},
				SelectedKey: "staging",
			})
		}))

		co.WithChild("settings", co.New(std.Button, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/settings.png"),
			})
		}))
	})
}
