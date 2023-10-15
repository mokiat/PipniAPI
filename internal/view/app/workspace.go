package app

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workflow"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	contextview "github.com/mokiat/PipniAPI/internal/view/context"
	endpointview "github.com/mokiat/PipniAPI/internal/view/endpoint"
	"github.com/mokiat/PipniAPI/internal/view/welcome"
	workflowview "github.com/mokiat/PipniAPI/internal/view/workflow"
	"github.com/mokiat/PipniAPI/internal/widget"
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
						Icon: c.editorImage(editor),
						Text: func() string {
							text := editor.Name()
							if editor.CanSave() {
								text += " *"
							}
							return text
						}(),
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

		switch editor := selectedEditor.(type) {
		case *context.Editor:
			co.WithChild(fmt.Sprintf("editor-%s", editor.ID()), co.New(contextview.Editor, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(contextview.EditorData{
					EditorModel: editor,
				})
			}))

		case *endpoint.Editor:
			co.WithChild(fmt.Sprintf("editor-%s", editor.ID()), co.New(endpointview.Editor, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(endpointview.EditorData{
					EditorModel: editor,
				})
			}))

		case *workflow.Editor:
			co.WithChild(fmt.Sprintf("editor-%s", editor.ID()), co.New(workflowview.Editor, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(workflowview.EditorData{
					EditorModel: editor,
				})
			}))

		case nil:
			co.WithChild("welcome-screen", co.New(welcome.Screen, func() {
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
	case workspace.EditorModifiedEvent:
		c.Invalidate()
	case registry.RegistryResourceRemovedEvent:
		c.closeEditorForResource(event.Resource)
	}
}

func (c *workspaceComponent) editorImage(editor workspace.Editor) *ui.Image {
	switch editor.(type) {
	case *context.Editor:
		return co.OpenImage(c.Scope(), "images/context.png")
	case *endpoint.Editor:
		return co.OpenImage(c.Scope(), "images/ping.png")
	case *workflow.Editor:
		return co.OpenImage(c.Scope(), "images/workflow.png")
	default:
		return nil
	}
}

func (c *workspaceComponent) selectEditor(editor workspace.Editor) {
	c.mdlWorkspace.SetSelectedID(editor.ID())
}

func (c *workspaceComponent) closeEditor(editor workspace.Editor, force bool) {
	if !force && editor.CanSave() {
		co.OpenOverlay(c.Scope(), co.New(widget.ConfirmationModal, func() {
			co.WithData(widget.ConfirmationModalData{
				Icon: co.OpenImage(c.Scope(), "images/warning.png"),
				Text: "There are unsaved changes!\n\nAre you sure you want to continue?",
			})
			co.WithCallbackData(widget.ConfirmationModalCallbackData{
				OnApply: func() {
					c.closeEditor(editor, true)
				},
			})
		}))
	} else {
		c.mdlWorkspace.RemoveEditor(editor)
	}
}

func (c *workspaceComponent) closeEditorForResource(resource registry.Resource) {
	for _, editor := range c.mdlWorkspace.Editors() {
		if editor.ID() == resource.ID() {
			c.closeEditor(editor, true)
		}
	}
}
