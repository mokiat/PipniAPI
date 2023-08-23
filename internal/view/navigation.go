package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewNavigationPanel(
	eventBus *mvc.EventBus,
	mdlWorkspace *model.Workspace,
) fyne.CanvasObject {

	registry := model.NewRegistry(eventBus) // TODO: Move outside

	top := NewEnvSelector(eventBus, registry, mdlWorkspace)
	middle := NewResourceTree(eventBus, registry)

	return container.NewBorder(top, nil, nil, nil, middle)
}
