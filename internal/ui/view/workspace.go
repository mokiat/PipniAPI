package view

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/PipniAPI/internal/model/workflow"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Workspace = mvc.EventListener(co.Define(&workspaceComponent{}))

type WorkspaceData struct {
	WorkspaceModel *workspace.Model
}

type workspaceComponent struct {
	co.BaseComponent

	mdlWorkspace *workspace.Model
}

func (c *workspaceComponent) OnUpsert() {
	data := co.GetData[WorkspaceData](c.Properties())
	c.mdlWorkspace = data.WorkspaceModel
}

func (c *workspaceComponent) Render() co.Instance {
	selectedEditor := c.mdlWorkspace.SelectedEditor()

	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			Layout: layout.Frame(),
		})

		co.WithChild("tabbar", co.New(std.Tabbar, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})

			c.mdlWorkspace.EachEditor(func(editor workspace.Editor) {
				co.WithChild(editor.ID(), co.New(std.TabbarTab, func() {
					co.WithData(std.TabbarTabData{
						Icon:     c.editorImage(editor),
						Text:     editor.Name(),
						Selected: c.mdlWorkspace.SelectedEditor() == editor,
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
			})
		}))

		if selectedEditor != nil {
			// TODO: Dynamic type based on workspace model editor selection
			co.WithChild(fmt.Sprintf("tabbar-editor-%s", selectedEditor.ID()), co.New(EndpointEditor, func() {
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
	case workspace.EditorAddedEvent:
		c.Invalidate()
	case workspace.EditorRemovedEvent:
		c.Invalidate()
	case workspace.EditorSelectedEvent:
		c.Invalidate()
	case registrymodel.RegistryResourceRemovedEvent:
		c.closeEditorForResource(event.Resource)
	}
}

func (c *workspaceComponent) editorImage(editor workspace.Editor) *ui.Image {
	switch editor.(type) {
	case *endpoint.Editor:
		return co.OpenImage(c.Scope(), "images/ping.png")
	case *workflow.Editor:
		return co.OpenImage(c.Scope(), "images/workflow.png")
	case *context.Editor:
		return co.OpenImage(c.Scope(), "images/settings.png")
	default:
		return nil
	}
}

func (c *workspaceComponent) selectEditor(editor workspace.Editor) {
	c.mdlWorkspace.SetSelectedID(editor.ID())
}

func (c *workspaceComponent) closeEditor(editor workspace.Editor, force bool) {
	// TODO: Check if dirty and open a confirmation dialog if dirty.
	c.mdlWorkspace.RemoveEditor(editor)
}

func (c *workspaceComponent) closeEditorForResource(resource registrymodel.Resource) {
	for _, editor := range c.mdlWorkspace.Editors() {
		if editor.ID() == resource.ID() {
			c.closeEditor(editor, true)
		}
	}
}
