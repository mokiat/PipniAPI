package context

import (
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	EditorID = "23dd8654-7af0-4a40-b63d-18211a7fa838"
)

func NewEditor(eventBus *mvc.EventBus, model *Model) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		model:    model,
	}
}

type Editor struct {
	workspace.NoSaveEditor

	eventBus *mvc.EventBus
	model    *Model
}

func (e *Editor) ID() string {
	return EditorID
}

func (e *Editor) Name() string {
	return "Context"
}
