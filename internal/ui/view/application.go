package view

import (
	"errors"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workflow"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
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
	mdlContext   *context.Model
	mdlRegistry  *registry.Model
	mdlWorkspace *workspace.Model
}

func (c *applicationComponent) OnCreate() {
	var loadErr error

	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	c.mdlContext = context.NewModel(c.eventBus, "context.json")
	if err := c.mdlContext.Load(); err != nil {
		c.mdlContext.Clear() // start with blank
		if !errors.Is(err, context.ErrContextNotFound) {
			loadErr = errors.Join(loadErr, err)
		}
	}

	c.mdlRegistry = registry.NewModel(c.eventBus, "registry.json")
	if err := c.mdlRegistry.Load(); err != nil {
		c.mdlRegistry.Clear() // start with blank
		if !errors.Is(err, registry.ErrRegistryNotFound) {
			loadErr = errors.Join(loadErr, err)
		}
	}

	c.mdlWorkspace = workspace.NewModel(c.eventBus)

	if loadErr != nil {
		log.Error("Error loading models: %v", loadErr)
		co.OpenOverlay(c.Scope(), co.New(NotificationModal, func() {
			co.WithData(NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: "The program encountered an error.\n\nSome of the state could not be restored.",
			})
		}))
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
						WorkspaceModel: c.mdlWorkspace,
						ContextModel:   c.mdlContext,
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
	case registry.RegistrySelectionChangedEvent:
		c.openEditorForRegistryItem(event.SelectedID)
	case workspace.EditorSelectedEvent:
		c.selectResourceForEditor(event.Editor)
	}
}

func (c *applicationComponent) openEditorForRegistryItem(itemID string) {
	if editor := c.mdlWorkspace.FindEditor(itemID); editor != nil {
		c.mdlWorkspace.SetSelectedID(editor.ID())
		return
	}

	resource := c.mdlRegistry.Root().FindResource(itemID)
	switch resource := resource.(type) {
	case *registry.Endpoint:
		c.mdlWorkspace.AppendEditor(endpoint.NewEditor(c.eventBus, resource))
	case *registry.Workflow:
		c.mdlWorkspace.AppendEditor(workflow.NewEditor(c.eventBus, resource))
	}
}

func (c *applicationComponent) selectResourceForEditor(editor workspace.Editor) {
	if editor != nil {
		c.mdlRegistry.SetSelectedID(editor.ID())
	} else {
		c.mdlRegistry.SetSelectedID("")
	}
}
