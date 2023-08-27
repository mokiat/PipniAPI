package view

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/ui/model"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Application = co.Define(&applicationComponent{})

type applicationComponent struct {
	co.BaseComponent

	mdlContext  *model.Context
	mdlRegistry *model.Registry
}

func (c *applicationComponent) OnCreate() {
	eventBus := co.TypedValue[*mvc.EventBus](c.Scope())

	c.mdlContext = model.NewContext(eventBus, "context.json")
	if err := c.mdlContext.Load(); err != nil {
		panic(fmt.Errorf("error loading registry: %w", err)) // TODO: Show error dialog and continue with blank state.
	}

	c.mdlRegistry = model.NewRegistry(eventBus, "registry.json")
	if err := c.mdlRegistry.Load(); err != nil {
		panic(fmt.Errorf("error loading registry: %w", err)) // TODO: Show error dialog and continue with blank state.
	}
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

				co.WithChild("environment-selection", co.New(EnvironmentSelection, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
					})
					co.WithData(EnvironmentSelectionData{
						ContextModel: c.mdlContext,
					})
				}))

				co.WithChild("endpoint-selection", co.New(EndpointSelection, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentCenter,
						VerticalAlignment:   layout.VerticalAlignmentCenter,
					})
					co.WithData(EndpointSelectionData{
						RegistryModel: c.mdlRegistry,
					})
				}))

				co.WithChild("endpoint-management", co.New(EndpointManagement, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentBottom,
					})
				}))
			}))

			co.WithChild("workspace", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ContainerData{
					Layout: layout.Frame(),
				})

				co.WithChild("tabbar", co.New(std.Tabbar, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
					})

					co.WithChild("tab-<id>", co.New(std.TabbarTab, func() {
						co.WithData(std.TabbarTabData{
							Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
							Text:     "List Users",
							Selected: true,
						})
					}))

					co.WithChild("tab-<id2>", co.New(std.TabbarTab, func() {
						co.WithData(std.TabbarTabData{
							Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
							Text:     "Get User",
							Selected: false,
						})
					}))
				}))

				// TODO: Dynamic based on workspace model editor selection
				co.WithChild("tabbar-editor-<id>", co.New(std.Container, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentCenter,
						VerticalAlignment:   layout.VerticalAlignmentCenter,
					})
					co.WithData(std.ContainerData{
						BackgroundColor: opt.V(ui.Purple()),
					})
				}))
			}))

		}))
	})
}
