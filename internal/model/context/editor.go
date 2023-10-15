package context

import (
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, reg *registry.Model, cont *registry.Context) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		cont:     cont,
	}
}

type Editor struct {
	workspace.NoSaveEditor

	eventBus *mvc.EventBus
	cont     *registry.Context
}

func (e *Editor) ID() string {
	return e.cont.ID()
}

func (e *Editor) Name() string {
	return e.cont.Name()
}
