package registry

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Item = co.Define(&itemComponent{})

type ItemData struct {
	Selected bool
	Icon     *ui.Image
	Text     string
}

type ItemCallbackData struct {
	OnClick std.OnActionFunc
}

type itemComponent struct {
	co.BaseComponent

	selected bool
	icon     *ui.Image
	text     string

	onClick std.OnActionFunc
}

func (c *itemComponent) OnUpsert() {
	data := co.GetData[ItemData](c.Properties())
	c.selected = data.Selected
	c.icon = data.Icon
	c.text = data.Text

	callbackData := co.GetCallbackData[ItemCallbackData](c.Properties())
	c.onClick = callbackData.OnClick
}

func (c *itemComponent) Render() co.Instance {
	return co.New(std.ListItem, func() {
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: c.onClick,
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
