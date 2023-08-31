package view

import (
	"errors"
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/PipniAPI/internal/ui/model"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Application = mvc.EventListener(co.Define(&applicationComponent{}))

type applicationComponent struct {
	co.BaseComponent

	eventBus     *mvc.EventBus
	mdlContext   *model.Context
	mdlRegistry  *registrymodel.Registry
	mdlWorkspace *model.Workspace
}

func (c *applicationComponent) OnCreate() {
	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	c.mdlContext = model.NewContext(c.eventBus, "context.json")
	if err := c.mdlContext.Load(); err != nil {
		panic(fmt.Errorf("error loading registry: %w", err)) // TODO: Show error dialog and continue with blank state.
	}

	c.mdlRegistry = registrymodel.NewRegistry(c.eventBus, "registry.json")
	if err := c.mdlRegistry.Load(); err != nil {
		c.mdlRegistry.Clear() // start with blank
		if !errors.Is(err, registrymodel.ErrRegistryNotFound) {
			panic("TODO: Show error message dialog")
		}
	}

	c.mdlWorkspace = model.NewWorkspace(c.eventBus)
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
					co.WithData(EndpointManagementData{
						RegistryModel: c.mdlRegistry,
					})
				}))
			}))

			co.WithChild("workspace", co.New(Workspace, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(WorkspaceData{
					WorkspaceModel: c.mdlWorkspace,
				})
			}))
		}))
	})
}

func (c *applicationComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case registrymodel.RegistrySelectionChangedEvent:
		c.openEditorForRegistryItem(event.SelectedID)
	case model.EditorSelectedEvent:
		c.selectResourceForEditor(event.Editor)
	case model.ContextEditorOpenEvent:
		c.openContextEditor()
	}
}

func (c *applicationComponent) openEditorForRegistryItem(itemID string) {
	resource := c.mdlRegistry.Root().FindResource(itemID)
	switch resource := resource.(type) {
	case *registrymodel.Endpoint:
		c.mdlWorkspace.OpenEditor(model.NewEndpointEditor(c.eventBus, resource))
	case *registrymodel.Workflow:
		c.mdlWorkspace.OpenEditor(model.NewWorkflowEditor(c.eventBus, resource))
	}
}

func (c *applicationComponent) openContextEditor() {
	c.mdlWorkspace.OpenEditor(model.NewContextEditor(c.eventBus, c.mdlContext))
}

func (c *applicationComponent) selectResourceForEditor(editor model.Editor) {
	if editor != nil {
		c.mdlRegistry.SetSelectedID(editor.ID())
	} else {
		c.mdlRegistry.SetSelectedID("")
	}
}
