package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"golang.org/x/exp/slices"
)

func NewWorkspace(eventBus *mvc.EventBus, mdl *model.Workspace) fyne.CanvasObject {
	contentForEditor := make(map[model.Editor]fyne.CanvasObject)

	tabs := container.NewDocTabs()

	updateTabs := func(editors []model.Editor) {
		items := make([]*container.TabItem, len(editors))
		for i, editor := range editors {
			content, ok := contentForEditor[editor]
			if !ok {
				content = createEditorContent(editor)
				contentForEditor[editor] = content
			}
			items[i] = container.NewTabItem(editor.Title(), content)
			// container.NewTabItemWithIcon(text string, icon fyne.Resource, content fyne.CanvasObject)
		}
		tabs.SetItems(items)
	}
	updateTabs(mdl.Editors())

	updateSelection := func(activeEditor model.Editor) {
		tabs.SelectIndex(slices.Index(mdl.Editors(), activeEditor))
	}
	updateSelection(mdl.ActiveEditor())

	tabs.CloseIntercept = func(ti *container.TabItem) {
		// TODO: If dirty, don't ask beforehand
		tabs.Remove(ti)
		tabs.OnClosed(ti)
	}

	tabs.OnClosed = func(ti *container.TabItem) {
		for editor, content := range contentForEditor {
			if content == ti.Content {
				mdl.RemoveEditor(editor)
			}
		}
	}

	eventBus.Subscribe(func(event mvc.Event) {
		switch event := event.(type) {
		case model.EditorAddedEvent:
			updateTabs(mdl.Editors())
		case model.EditorRemovedEvent:
			updateTabs(mdl.Editors())
		case model.EditorSelectedEvent:
			updateSelection(event.Editor)
		}
	})

	// tabs.SetItems(items []*container.TabItem)

	// return container.NewAppTabs(
	// 	container.NewTabItem("hello", container.NewVSplit(
	// 		widget.NewLabel("request"),
	// 		widget.NewLabel("response"),
	// 	)),
	// 	container.NewTabItem("bye", container.NewVSplit(
	// 		widget.NewLabel("env"),
	// 		widget.NewLabel("tmp"),
	// 	)),
	// )
	return tabs
}

func createEditorContent(editor model.Editor) fyne.CanvasObject {
	switch editor.(type) {
	case *model.EnvironmentEditor:
		return widget.NewLabel("Environment Configuration ...")
	case *model.RequestEditor:
		return widget.NewLabel("Request / Response ...")
	default:
		return widget.NewLabel("Unsupported editor")
	}
}
