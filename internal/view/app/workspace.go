package app

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workflow"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/PipniAPI/internal/shortcuts"
	contextview "github.com/mokiat/PipniAPI/internal/view/context"
	endpointview "github.com/mokiat/PipniAPI/internal/view/endpoint"
	"github.com/mokiat/PipniAPI/internal/view/welcome"
	workflowview "github.com/mokiat/PipniAPI/internal/view/workflow"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/debug/log"
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

var _ ui.ElementKeyboardHandler = (*workspaceComponent)(nil)
var _ ui.ElementStateHandler = (*workspaceComponent)(nil)

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

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    layout.Frame(),
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

func (c *workspaceComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	os := element.Window().Platform().OS()
	if shortcuts.IsClose(os, event) {
		if event.Action == ui.KeyboardActionDown {
			if selectedEditor := c.mdlWorkspace.SelectedEditor(); selectedEditor != nil {
				c.closeEditor(selectedEditor, false)
			} else if len(c.mdlWorkspace.Editors()) == 0 {
				co.Window(c.Scope()).Close()
			}
		}
		return true
	}
	if shortcuts.IsSave(os, event) {
		if event.Action == ui.KeyboardActionDown {
			if selectedEditor := c.mdlWorkspace.SelectedEditor(); selectedEditor != nil {
				c.saveEditor(selectedEditor)
			}
		}
	}
	return false
}

func (c *workspaceComponent) OnSave(element *ui.Element) bool {
	if selectedEditor := c.mdlWorkspace.SelectedEditor(); selectedEditor != nil {
		c.saveEditor(selectedEditor)
		return true
	}
	return false
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

func (c *workspaceComponent) saveEditor(editor workspace.Editor) {
	if err := editor.Save(); err != nil {
		log.Error("Error saving editor changes: %v", err)
		co.OpenOverlay(c.Scope(), co.New(widget.NotificationModal, func() {
			co.WithData(widget.NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: "The program encountered an error.\n\nChanges could not be saved.",
			})
		}))
	}
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
