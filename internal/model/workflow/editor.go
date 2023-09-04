package workflow

import (
	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, workflow *registrymodel.Workflow) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		workflow: workflow,
	}
}

type Editor struct {
	workspace.NoSaveEditor
	workspace.NoHistoryEditor

	eventBus *mvc.EventBus
	workflow *registrymodel.Workflow
}

func (e *Editor) ID() string {
	return e.workflow.ID()
}

func (e *Editor) Name() string {
	return e.workflow.Name()
}
