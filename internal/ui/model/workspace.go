package model

import (
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui/mvc"
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

func (w *Workspace) OpenEditor(editor Editor) {
	existingEditor, ok := gog.FindFunc(w.editors, func(candidate Editor) bool {
		return candidate.ID() == editor.ID()
	})
	if !ok {
		w.editors = append(w.editors, editor)
		existingEditor = editor
		w.eventBus.Notify(EditorAddedEvent{
			Editor: editor,
		})
	}
	w.SetActiveEditor(existingEditor)
}

func (w *Workspace) CloseEditor(editor Editor) {
	index := slices.Index(w.editors, editor)
	if index < 0 {
		return
	}

	if editor == w.activeEditor {
		if index > 0 {
			w.SetActiveEditor(w.editors[index-1])
		} else if len(w.editors) > 1 {
			w.SetActiveEditor(w.editors[1])
		} else {
			w.SetActiveEditor(nil)
		}
	}

	w.editors = slices.Delete(w.editors, index, index+1)
	w.eventBus.Notify(EditorRemovedEvent{
		Editor: editor,
	})
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
