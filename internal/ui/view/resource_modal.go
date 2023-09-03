package view

import (
	"strings"

	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var ResourceModal = co.Define(&resourceModalComponent{})

type ResourceModalData struct {
	Name          string
	Kind          registrymodel.ResourceKind
	CanChangeKind bool
}

type ResourceModalCallbackData struct {
	OnApply func(name string, kind registrymodel.ResourceKind)
}

type resourceModalComponent struct {
	co.BaseComponent

	name          string
	kind          registrymodel.ResourceKind
	canChangeKind bool

	onApply func(name string, kind registrymodel.ResourceKind)
}

func (c *resourceModalComponent) OnCreate() {
	data := co.GetData[ResourceModalData](c.Properties())
	c.name = data.Name
	c.kind = data.Kind
	c.canChangeKind = data.CanChangeKind

	callbackData := co.GetCallbackData[ResourceModalCallbackData](c.Properties())
	c.onApply = callbackData.OnApply
}

func (c *resourceModalComponent) Render() co.Instance {
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
						Selected: c.kind == registrymodel.ResourceKindEndpoint,
						Enabled:  opt.V(c.canChangeKind),
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setKind(registrymodel.ResourceKindEndpoint)
						},
					})
				}))

				co.WithChild("workflow", co.New(std.ToolbarButton, func() {
					co.WithData(std.ToolbarButtonData{
						Icon:     co.OpenImage(c.Scope(), "images/workflow.png"),
						Text:     "Workflow",
						Selected: c.kind == registrymodel.ResourceKindWorkflow,
						Enabled:  opt.V(c.canChangeKind),
					})
					co.WithCallbackData(std.ToolbarButtonCallbackData{
						OnClick: func() {
							c.setKind(registrymodel.ResourceKindWorkflow)
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

					co.WithChild("editbox", co.New(std.Editbox, func() {
						co.WithLayoutData(layout.Data{
							HorizontalAlignment: layout.HorizontalAlignmentCenter,
						})
						co.WithData(std.EditboxData{
							Text: c.name,
						})
						co.WithCallbackData(std.EditboxCallbackData{
							OnChanged: func(text string) {
								c.setName(text)
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

func (c *resourceModalComponent) resourceInfo(kind registrymodel.ResourceKind) string {
	// TODO: Manual text wrapping won't be needed here once a TextArea is used.
	switch kind {
	case registrymodel.ResourceKindEndpoint:
		return "An Endpoint resource represents a specific RESTful endpoint\ninvocation.\n\nEnter a name to properly reflect the endpoint's purpose."
	case registrymodel.ResourceKindWorkflow:
		return "A Workflow can be used to orchestrate the invocation of a\nnumber of Endpoint resources by passing data between them\nand controlling the call order.\n\nEnter a name to properly reflect the workflow's purpose."
	default:
		return ""
	}
}

func (c *resourceModalComponent) setName(name string) {
	c.name = name
	c.Invalidate()
}

func (c *resourceModalComponent) setKind(kind registrymodel.ResourceKind) {
	c.kind = kind
	c.Invalidate()
}

func (c *resourceModalComponent) onGo() {
	c.onApply(c.name, c.kind)
	co.CloseOverlay(c.Scope())
}

func (c *resourceModalComponent) onCancel() {
	co.CloseOverlay(c.Scope())
}
