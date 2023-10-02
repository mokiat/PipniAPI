package workflow

import (
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, workflow *registry.Workflow) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		workflow: workflow,
	}
}

type Editor struct {
	workspace.NoSaveEditor

	eventBus *mvc.EventBus
	workflow *registry.Workflow
}

func (e *Editor) ID() string {
	return e.workflow.ID()
}

func (e *Editor) Name() string {
	return e.workflow.Name()
}
