package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewWindow(eventBus *mvc.EventBus) fyne.CanvasObject {
	return container.NewHSplit(
		NewTreeMenu(eventBus),
		NewContentArena(eventBus),
	)
}

func NewTreeMenu(eventBus *mvc.EventBus) fyne.CanvasObject {
	registry := model.NewRegistry(eventBus) // TODO: Move outside

	top := NewEnvSelector(eventBus, registry)
	middle := NewResourceTree(eventBus, registry)

	return container.NewBorder(top, nil, nil, nil, middle)
}

func NewContentArena(eventBus *mvc.EventBus) fyne.CanvasObject {
	return container.NewVSplit(
		widget.NewLabel("request"),
		widget.NewLabel("response"),
	)
}
