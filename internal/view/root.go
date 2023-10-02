package view

import (
	"errors"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workflow"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	appview "github.com/mokiat/PipniAPI/internal/view/app"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Root = mvc.EventListener(co.Define(&rootComponent{}))

type rootComponent struct {
	co.BaseComponent

	eventBus     *mvc.EventBus
	mdlContext   *context.Model
	mdlRegistry  *registry.Model
	mdlWorkspace *workspace.Model
}

func (c *rootComponent) OnCreate() {
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
		co.OpenOverlay(c.Scope(), co.New(widget.NotificationModal, func() {
			co.WithData(widget.NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: "The program encountered an error.\n\nSome of the state could not be restored.",
			})
		}))
	}
}

func (c *rootComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Frame(),
		})

		co.WithChild("toolbar", co.New(appview.Toolbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(appview.ToolbarData{
				WorkspaceModel: c.mdlWorkspace,
			})
		}))

		co.WithChild("content", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ElementData{
				Layout: layout.Frame(),
			})

			co.WithChild("drawer", co.New(appview.Drawer, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(300),
				})
				co.WithData(appview.DrawerData{
					ContextModel:   c.mdlContext,
					RegistryModel:  c.mdlRegistry,
					WorkspaceModel: c.mdlWorkspace,
				})
			}))

			co.WithChild("workspace", co.New(appview.Workspace, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(appview.WorkspaceData{
					WorkspaceModel: c.mdlWorkspace,
					ContextModel:   c.mdlContext,
				})
			}))
		}))
	})
}

func (c *rootComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case registry.RegistrySelectionChangedEvent:
		c.openEditorForRegistryItem(event.SelectedID)
	case workspace.EditorSelectedEvent:
		c.selectResourceForEditor(event.Editor)
	}
}

func (c *rootComponent) openEditorForRegistryItem(itemID string) {
	if editor := c.mdlWorkspace.FindEditor(itemID); editor != nil {
		c.mdlWorkspace.SetSelectedID(editor.ID())
		return
	}

	resource := c.mdlRegistry.Root().FindResource(itemID)
	switch resource := resource.(type) {
	case *registry.Endpoint:
		c.mdlWorkspace.AppendEditor(endpoint.NewEditor(c.eventBus, c.mdlRegistry, resource))
	case *registry.Workflow:
		c.mdlWorkspace.AppendEditor(workflow.NewEditor(c.eventBus, resource))
	}
}

func (c *rootComponent) selectResourceForEditor(editor workspace.Editor) {
	if editor != nil {
		c.mdlRegistry.SetSelectedID(editor.ID())
	} else {
		c.mdlRegistry.SetSelectedID("")
	}
}
