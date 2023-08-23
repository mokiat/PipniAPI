package model

import (
	"github.com/mokiat/PipniAPI/internal/mvc"
	"golang.org/x/exp/slices"
)

func NewWorkspace(eventBus *mvc.EventBus) *Workspace {
	return &Workspace{
		eventBus: eventBus,
	}
}

type Workspace struct {
	eventBus *mvc.EventBus

	editors      []Editor
	activeEditor Editor
}

func (w *Workspace) Editors() []Editor {
	return w.editors
}

func (w *Workspace) AddEditor(editor Editor) {
	w.editors = append(w.editors, editor)
	w.eventBus.Notify(EditorAddedEvent{
		Editor: editor,
	})
	if len(w.editors) == 1 {
		w.SetActiveEditor(editor)
	}
}

func (w *Workspace) RemoveEditor(editor Editor) {
	w.editors = slices.DeleteFunc(w.editors, func(candidate Editor) bool {
		return candidate == editor
	})
	w.eventBus.Notify(EditorRemovedEvent{
		Editor: editor,
	})
	if w.activeEditor == editor {
		if len(w.editors) > 0 {
			w.activeEditor = w.editors[0]
		} else {
			w.activeEditor = nil
		}
	}
}

func (w *Workspace) ActiveEditor() Editor {
	return w.activeEditor
}

func (w *Workspace) SetActiveEditor(editor Editor) {
	w.activeEditor = editor
	w.eventBus.Notify(EditorSelectedEvent{
		Editor: editor,
	})
}

type EditorAddedEvent struct {
	Editor Editor
}

type EditorRemovedEvent struct {
	Editor Editor
}

type EditorSelectedEvent struct {
	Editor Editor
}
