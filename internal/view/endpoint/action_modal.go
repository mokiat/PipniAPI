package endpoint

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var ActionModal = co.Define(&actionModalComponent{})

type ActionModalCallbackData struct {
	OnCancel std.OnActionFunc
}

type actionModalComponent struct {
	co.BaseComponent

	onCancel std.OnActionFunc
}

func (c *actionModalComponent) OnCreate() {
	callbackData := co.GetCallbackData[ActionModalCallbackData](c.Properties())
	c.onCancel = callbackData.OnCancel
}

func (c *actionModalComponent) Render() co.Instance {
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
						Image:      co.OpenImage(c.Scope(), "images/info.png"),
						ImageColor: opt.V(ui.Black()),
						Mode:       std.ImageModeFit,
					})
				}))

				co.WithChild("text", co.New(std.Label, func() {
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(20)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      "Request in progress...\n\nClick Cancel to stop.",
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

				co.WithChild("cancel", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text: "Cancel",
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onCancel,
					})
				}))
			}))
		}))
	})
}
