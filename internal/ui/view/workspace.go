package view

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/PipniAPI/internal/ui/model"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Workspace = mvc.EventListener(co.Define(&workspaceComponent{}))

type WorkspaceData struct {
	WorkspaceModel *model.Workspace
}

type workspaceComponent struct {
	co.BaseComponent

	mdlWorkspace *model.Workspace
}

func (c *workspaceComponent) OnUpsert() {
	data := co.GetData[WorkspaceData](c.Properties())
	c.mdlWorkspace = data.WorkspaceModel
}

func (c *workspaceComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			Layout: layout.Frame(),
		})

		co.WithChild("tabbar", co.New(std.Tabbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			for _, editor := range c.mdlWorkspace.Editors() {
				editor := editor

				co.WithChild(editor.ID(), co.New(std.TabbarTab, func() {
					co.WithData(std.TabbarTabData{
						Icon:     c.editorImage(editor),
						Text:     editor.Title(),
						Selected: c.mdlWorkspace.ActiveEditor() == editor,
					})
					co.WithCallbackData(std.TabbarTabCallbackData{
						OnClick: func() {
							c.selectEditor(editor)
						},
						OnClose: func() {
							c.closeEditor(editor, false)
						},
					})
				}))
			}
		}))

		if activeEditor := c.mdlWorkspace.ActiveEditor(); activeEditor != nil {
			// TODO: Dynamic type based on workspace model editor selection
			co.WithChild(fmt.Sprintf("tabbar-editor-%s", activeEditor.ID()), co.New(EndpointEditor, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
			}))
		} else {
			co.WithChild("welcome-screen", co.New(WelcomeScreen, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
			}))
		}
	})
}

func (c *workspaceComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case model.EditorSelectedEvent:
		c.Invalidate()
	case model.EditorAddedEvent:
		c.Invalidate()
	case model.EditorRemovedEvent:
		c.Invalidate()
	case registrymodel.RegistryResourceRemovedEvent:
		c.closeEditorForResource(event.Resource)
	}
}

func (c *workspaceComponent) editorImage(editor model.Editor) *ui.Image {
	switch editor.(type) {
	case *model.EndpointEditor:
		return co.OpenImage(c.Scope(), "images/ping.png")
	case *model.WorkflowEditor:
		return co.OpenImage(c.Scope(), "images/workflow.png")
	case *model.ContextEditor:
		return co.OpenImage(c.Scope(), "images/settings.png")
	default:
		return nil
	}
}

func (c *workspaceComponent) selectEditor(editor model.Editor) {
	c.mdlWorkspace.SetActiveEditor(editor)
}

func (c *workspaceComponent) closeEditor(editor model.Editor, force bool) {
	// TODO: Check if dirty and open a confirmation dialog if dirty.
	c.mdlWorkspace.CloseEditor(editor)
}

func (c *workspaceComponent) closeEditorForResource(resource registrymodel.Resource) {
	for _, editor := range c.mdlWorkspace.Editors() {
		if editor.ID() == resource.ID() {
			c.closeEditor(editor, true)
		}
	}
}
