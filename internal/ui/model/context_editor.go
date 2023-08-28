package model

import "github.com/mokiat/lacking/ui/mvc"

func NewContextEditor(eventBus *mvc.EventBus, context *Context) Editor {
	return &ContextEditor{
		eventBus: eventBus,
		context:  context,
	}
}

type ContextEditor struct {
	eventBus *mvc.EventBus
	context  *Context
}

func (e *ContextEditor) ID() string {
	return "-context-"
}

func (e *ContextEditor) Title() string {
	return "Environments"
}

type ContextEditorOpenEvent struct {
	Context *Context
}
