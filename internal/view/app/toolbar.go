package app

import (
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/PipniAPI/internal/view/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = mvc.EventListener(co.Define(&toolbarComponent{}))

type ToolbarData struct {
	WorkspaceModel *workspace.Model
}

type toolbarComponent struct {
	co.BaseComponent

	mdlWorkspace *workspace.Model
}

func (c *toolbarComponent) OnUpsert() {
	data := co.GetData[ToolbarData](c.Properties())
	c.mdlWorkspace = data.WorkspaceModel
}

func (c *toolbarComponent) Render() co.Instance {
	editor := c.mdlWorkspace.SelectedEditor()
	canSave := (editor != nil) && (editor.CanSave())
	canUndo := (editor != nil) && (editor.CanUndo())
	canRedo := (editor != nil) && (editor.CanRedo())

	return co.New(std.Toolbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("logo", co.New(std.ToolbarLogo, func() {
			co.WithData(std.ToolbarLogoData{
				Image: co.OpenImage(c.Scope(), "images/icon.png"),
				Text:  "Pipni API",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.importToRegistry,
			})
		}))

		co.WithChild("separator-after-logo", co.New(std.ToolbarSeparator, nil))

		co.WithChild("import", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/import.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.importToRegistry,
			})
		}))

		co.WithChild("export", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/export.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.exportFromRegistry,
			})
		}))

		co.WithChild("separator-after-import-export", co.New(std.ToolbarSeparator, nil))

		co.WithChild("cut", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/cut.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.cutContent,
			})
		}))

		co.WithChild("copy", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/copy.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.copyContent,
			})
		}))

		co.WithChild("paste", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/paste.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.pasteContent,
			})
		}))

		co.WithChild("separator-after-copy-paste", co.New(std.ToolbarSeparator, nil))

		co.WithChild("save", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/save.png"),
				Enabled: opt.V(canSave),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.saveEditorChanges(editor)
				},
			})
		}))

		co.WithChild("separator-after-save", co.New(std.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/undo.png"),
				Enabled: opt.V(canUndo),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.undoEditorChange(editor)
				},
			})
		}))

		co.WithChild("redo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/redo.png"),
				Enabled: opt.V(canRedo),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: func() {
					c.redoEditorChange(editor)
				},
			})
		}))
	})
}

func (c *toolbarComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case workspace.EditorSelectedEvent:
		c.Invalidate()
	case workspace.EditorModifiedEvent:
		c.Invalidate()
	}
}

func (c *toolbarComponent) importToRegistry() {
	log.Info("Import")
}

func (c *toolbarComponent) exportFromRegistry() {
	log.Info("Export")
}

func (c *toolbarComponent) cutContent() {
	log.Info("Cut")
}

func (c *toolbarComponent) copyContent() {
	log.Info("Copy")
}

func (c *toolbarComponent) pasteContent() {
	log.Info("Paste")
}

func (c *toolbarComponent) saveEditorChanges(editor workspace.Editor) {
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

func (c *toolbarComponent) undoEditorChange(editor workspace.Editor) {
	editor.Undo()
}

func (c *toolbarComponent) redoEditorChange(editor workspace.Editor) {
	editor.Redo()
}
