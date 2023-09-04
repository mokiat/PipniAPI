package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var NotificationModal = co.Define(&notificationModalComponent{})

type NotificationModalData struct {
	Icon *ui.Image
	Text string
}

type notificationModalComponent struct {
	co.BaseComponent

	icon *ui.Image
	text string
}

func (c *notificationModalComponent) OnCreate() {
	data := co.GetData[NotificationModalData](c.Properties())
	c.icon = data.Icon
	c.text = data.Text
}

func (c *notificationModalComponent) Render() co.Instance {
	return co.New(std.Modal, func() {
		co.WithLayoutData(layout.Data{
			Width:            opt.V(500),
			Height:           opt.V(400),
			HorizontalCenter: opt.V(0),
			VerticalCenter:   opt.V(0),
		})

		co.WithChild("dialog", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Frame(layout.FrameSettings{
					ContentSpacing: ui.SymmetricSpacing(0, 20),
				}),
			})

			co.WithChild("content", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ElementData{
					Layout: layout.Frame(layout.FrameSettings{
						ContentSpacing: ui.Spacing{
							Left: 5,
						},
					}),
				})

				co.WithChild("icon", co.New(std.Picture, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
						Width:             opt.V(48),
						Height:            opt.V(48),
					})
					co.WithData(std.PictureData{
						Image:      c.icon,
						ImageColor: opt.V(ui.Black()),
						Mode:       std.ImageModeFit,
					})
				}))

				co.WithChild("text", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(20)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      c.text,
					})
				}))
			}))

			co.WithChild("footer", co.New(std.Toolbar, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentBottom,
				})
				co.WithData(std.ToolbarData{
					Positioning: std.ToolbarPositioningBottom,
				})

				co.WithChild("close", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text: "Close",
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onClose,
					})
				}))
			}))
		}))
	})
}

func (c *notificationModalComponent) onClose() {
	co.CloseOverlay(c.Scope())
}
