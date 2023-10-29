package registry

import (
	"strings"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Modal = co.Define(&modalComponent{})

type ModalData struct {
	Name          string
	Kind          registry.ResourceKind
	CanChangeKind bool
}

type ModalCallbackData struct {
	OnApply func(name string, kind registry.ResourceKind)
}

type modalComponent struct {
	co.BaseComponent

	name          string
	kind          registry.ResourceKind
	canChangeKind bool

	onApply func(name string, kind registry.ResourceKind)
}

func (c *modalComponent) OnCreate() {
	data := co.GetData[ModalData](c.Properties())
	c.name = data.Name
	c.kind = data.Kind
	c.canChangeKind = data.CanChangeKind

	callbackData := co.GetCallbackData[ModalCallbackData](c.Properties())
	c.onApply = callbackData.OnApply
}

func (c *modalComponent) Render() co.Instance {
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

			co.WithChild("header", co.New(std.Toolbar, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentTop,
				})

				co.WithChild("endpoint", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
						Text:     "Endpoint",
						Selected: c.kind == registry.ResourceKindEndpoint,
						Enabled:  opt.V(c.canChangeKind),
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setKind(registry.ResourceKindEndpoint)
						},
					})
				}))

				co.WithChild("workflow", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope(), "images/workflow.png"),
						Text:     "Workflow",
						Selected: c.kind == registry.ResourceKindWorkflow,
						Enabled:  opt.V(c.canChangeKind),
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setKind(registry.ResourceKindWorkflow)
						},
					})
				}))

				co.WithChild("context", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope(), "images/context.png"),
						Text:     "Context",
						Selected: c.kind == registry.ResourceKindContext,
						Enabled:  opt.V(c.canChangeKind),
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setKind(registry.ResourceKindContext)
						},
					})
				}))
			}))

			co.WithChild("content", co.New(std.Element, func() {
				co.WithLayoutData(layout.Data{
					VerticalAlignment: layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ElementData{
					Layout: layout.Vertical(layout.VerticalSettings{
						ContentSpacing: 30,
					}),
				})

				co.WithChild("settings", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Padding: ui.UniformSpacing(10),
						Layout: layout.Frame(layout.FrameSettings{
							ContentSpacing: ui.Spacing{
								Left: 10,
							},
						}),
					})

					co.WithChild("label", co.New(std.Label, func() {
						co.WithLayoutData(layout.Data{
							HorizontalAlignment: layout.HorizontalAlignmentLeft,
						})
						co.WithData(std.LabelData{
							Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
							FontSize:  opt.V(float32(18)),
							FontColor: opt.V(std.OnSurfaceColor),
							Text:      "Name:",
						})
					}))

					co.WithChild("editbox", co.New(widget.EditBox, func() {
						co.WithLayoutData(layout.Data{
							HorizontalAlignment: layout.HorizontalAlignmentCenter,
						})
						co.WithData(widget.EditBoxData{
							Text: c.name,
						})
						co.WithCallbackData(widget.EditBoxCallbackData{
							OnChange: func(text string) {
								c.setName(text)
							},
							OnSubmit: func(text string) {
								c.setName(text)
								c.onGo()
							},
						})
					}))
				}))

				co.WithChild("info", co.New(std.Label, func() {
					// TODO: Use TextArea component that automatically handles
					// text wrapping.
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(18)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      c.resourceInfo(c.kind),
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

				co.WithChild("go", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Text:    "Go",
						Enabled: opt.V(strings.TrimSpace(c.name) != ""),
					})
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentRight,
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: c.onGo,
					})
				}))

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

func (c *modalComponent) resourceInfo(kind registry.ResourceKind) string {
	// TODO: Manual text wrapping won't be needed here once a TextArea is used.
	switch kind {
	case registry.ResourceKindEndpoint:
		return "An Endpoint resource represents a specific RESTful endpoint\ninvocation.\n\nEnter a name to properly reflect the endpoint's purpose."
	case registry.ResourceKindWorkflow:
		return "A Workflow can be used to orchestrate the invocation of a\nnumber of Endpoint resources by passing data between them\nand controlling the call order.\n\nEnter a name to properly reflect the workflow's purpose."
	case registry.ResourceKindContext:
		return "A Context can be used to specify reusable parameters."
	default:
		return ""
	}
}

func (c *modalComponent) setName(name string) {
	c.name = name
	c.Invalidate()
}

func (c *modalComponent) setKind(kind registry.ResourceKind) {
	c.kind = kind
	c.Invalidate()
}

func (c *modalComponent) onGo() {
	c.onApply(c.name, c.kind)
	co.CloseOverlay(c.Scope())
}

func (c *modalComponent) onCancel() {
	co.CloseOverlay(c.Scope())
}
