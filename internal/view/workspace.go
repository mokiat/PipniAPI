package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"golang.org/x/exp/slices"
)

func (w *Window) newWorkspace() fyne.CanvasObject {
	contentForEditor := make(map[model.Editor]fyne.CanvasObject)

	tabs := container.NewDocTabs()

	updateTabs := func(editors []model.Editor) {
		items := make([]*container.TabItem, len(editors))
		for i, editor := range editors {
			content, ok := contentForEditor[editor]
			if !ok {
				content = w.createEditorContent(editor)
				contentForEditor[editor] = content
			}
			items[i] = container.NewTabItem(editor.Title(), content)
			// container.NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject)
		}
		tabs.SetItems(items)
	}
	updateTabs(w.mdlWorkspace.Editors())

	updateSelection := func(activeEditor model.Editor) {
		tabs.SelectIndex(slices.Index(w.mdlWorkspace.Editors(), activeEditor))
	}
	updateSelection(w.mdlWorkspace.ActiveEditor())

	tabs.CloseIntercept = func(ti *container.TabItem) {
		// TODO: If dirty, don't ask beforehand
		tabs.Remove(ti)
		tabs.OnClosed(ti)
	}

	tabs.OnClosed = func(ti *container.TabItem) {
		for editor, content := range contentForEditor {
			if content == ti.Content {
				w.mdlWorkspace.RemoveEditor(editor)
			}
		}
	}

	w.eventBus.Subscribe(func(event mvc.Event) {
		switch event := event.(type) {
		case model.EditorAddedEvent:
			updateTabs(w.mdlWorkspace.Editors())
		case model.EditorRemovedEvent:
			updateTabs(w.mdlWorkspace.Editors())
		case model.EditorSelectedEvent:
			updateSelection(event.Editor)
		}
	})

	return tabs
}

func (w *Window) createEditorContent(editor model.Editor) fyne.CanvasObject {
	switch editor := editor.(type) {
	case *model.EnvironmentEditor:
		return widget.NewLabel("Environment Configuration ...")
	case *model.EndpointEditor:
		return w.newEndpointEditor(editor)
	default:
		return widget.NewLabel("Unsupported editor")
	}
}
