package view

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = co.Define(&toolbarComponent{})

type ToolbarData struct{}

type toolbarComponent struct {
	co.BaseComponent
}

func (c *toolbarComponent) Render() co.Instance {
	return co.New(std.Toolbar, func() {
		co.WithLayoutData(c.Properties().LayoutData())

		co.WithChild("logo", co.New(std.ToolbarLogo, func() {
			co.WithData(std.ToolbarLogoData{
				Image: co.OpenImage(c.Scope(), "images/icon.png"),
				Text:  "Pipni API",
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onImport,
			})
		}))

		co.WithChild("separator-after-logo", co.New(std.ToolbarSeparator, nil))

		co.WithChild("import", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/import.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onImport,
			})
		}))

		co.WithChild("export", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/export.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onExport,
			})
		}))

		co.WithChild("separator-after-import-export", co.New(std.ToolbarSeparator, nil))

		co.WithChild("cut", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/cut.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onCut,
			})
		}))

		co.WithChild("copy", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/copy.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onCopy,
			})
		}))

		co.WithChild("paste", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/paste.png"),
				Enabled: opt.V(false),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onPaste,
			})
		}))

		co.WithChild("separator-after-copy-paste", co.New(std.ToolbarSeparator, nil))

		co.WithChild("save", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/save.png"),
				Enabled: opt.V(false),
				// Enabled: opt.V(c.history.CanSave()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onSave,
			})
		}))

		co.WithChild("separator-after-save", co.New(std.ToolbarSeparator, nil))

		co.WithChild("undo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/undo.png"),
				Enabled: opt.V(false),
				// Enabled: opt.V(c.history.CanUndo()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onUndo,
			})
		}))

		co.WithChild("redo", co.New(std.ToolbarButton, func() {
			co.WithData(std.ToolbarButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/redo.png"),
				Enabled: opt.V(false),
				// Enabled: opt.V(c.history.CanRedo()),
			})
			co.WithCallbackData(std.ToolbarButtonCallbackData{
				OnClick: c.onRedo,
			})
		}))
	})
}

func (c *toolbarComponent) onImport() {
	log.Info("Import")
}

func (c *toolbarComponent) onExport() {
	log.Info("Export")
}

func (c *toolbarComponent) onCut() {
	log.Info("Cut")
}

func (c *toolbarComponent) onCopy() {
	log.Info("Copy")
}

func (c *toolbarComponent) onPaste() {
	log.Info("Paste")
}

func (c *toolbarComponent) onSave() {
	log.Info("Save")
}

func (c *toolbarComponent) onUndo() {
	log.Info("Undo")
}

func (c *toolbarComponent) onRedo() {
	log.Info("Redo")
}
