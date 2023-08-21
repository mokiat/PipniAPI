package model

import "github.com/mokiat/PipniAPI/internal/mvc"

type Window struct {
	eventBus *mvc.EventBus
}

func (w *Window) Editors() []*Editor {
	return nil
}
