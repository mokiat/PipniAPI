package model

import (
	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewWorkflowEditor(eventBus *mvc.EventBus, workflow *registrymodel.Workflow) Editor {
	return &WorkflowEditor{
		eventBus: eventBus,
		workflow: workflow,
	}
}

type WorkflowEditor struct {
	eventBus *mvc.EventBus
	workflow *registrymodel.Workflow
}

func (e *WorkflowEditor) ID() string {
	return e.workflow.ID()
}

func (e *WorkflowEditor) Title() string {
	return e.workflow.Name()
}
