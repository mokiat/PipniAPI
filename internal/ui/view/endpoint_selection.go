package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var EndpointSelection = co.Define(&endpointSelectionComponent{})

type endpointSelectionComponent struct {
	co.BaseComponent
}

func (c *endpointSelectionComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.SymmetricSpacing(0, 1),
			BorderColor: opt.V(std.OutlineColor),
			Layout:      layout.Fill(),
			Padding:     ui.SymmetricSpacing(0, 5),
		})

		co.WithChild("list", co.New(std.List, func() {

			co.WithChild("item 01", co.New(EndpointItem, func() {
				co.WithData(EndpointItemData{
					Selected: true,
					Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
					Text:     "Create User",
				})
			}))

			co.WithChild("item 02", co.New(EndpointItem, func() {
				co.WithData(EndpointItemData{
					Selected: false,
					Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
					Text:     "Update User",
				})
			}))

			co.WithChild("item 03", co.New(EndpointItem, func() {
				co.WithData(EndpointItemData{
					Selected: false,
					Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
					Text:     "Delete User",
				})
			}))

			co.WithChild("item 04", co.New(EndpointItem, func() {
				co.WithData(EndpointItemData{
					Selected: false,
					Icon:     co.OpenImage(c.Scope(), "images/workflow.png"),
					Text:     "User Test Scenario",
				})
			}))
		}))
	})
}

var EndpointItem = co.Define(&endpointItemComponent{})

type EndpointItemData struct {
	Selected bool
	Icon     *ui.Image
	Text     string
}

type endpointItemComponent struct {
	co.BaseComponent

	selected bool
	icon     *ui.Image
	text     string
}

func (c *endpointItemComponent) OnUpsert() {
	data := co.GetData[EndpointItemData](c.Properties())
	c.selected = data.Selected
	c.icon = data.Icon
	c.text = data.Text
}

func (c *endpointItemComponent) Render() co.Instance {
	return co.New(std.ListItem, func() {
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})

		co.WithChild("holder", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Frame(),
			})

			if c.icon != nil {
				co.WithChild("icon", co.New(std.Picture, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentLeft,
						VerticalAlignment:   layout.VerticalAlignmentCenter,
						Width:               opt.V(24),
						Height:              opt.V(24),
					})
					co.WithData(std.PictureData{
						Image:      c.icon,
						ImageColor: opt.V(std.OnSurfaceColor),
						Mode:       std.ImageModeFit,
					})
				}))
			}

			if c.text != "" {
				co.WithChild("label", co.New(std.Label, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentCenter,
					})
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(20.0)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      c.text,
					})
				}))
			}
		}))
	})
}
