package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewWindow(eventBus *mvc.EventBus) fyne.CanvasObject {
	mdlWorkspace := model.NewWorkspace(eventBus)
	mdlWorkspace.AddEditor(model.NewEnvironmentEditor())
	mdlWorkspace.AddEditor(model.NewRequestEditor())

	return container.NewHSplit(
		NewNavigationPanel(eventBus, mdlWorkspace),
		NewWorkspace(eventBus, mdlWorkspace),
	)
}

// func NewTreeMenu(eventBus *mvc.EventBus) fyne.CanvasObject {
// 	registry := model.NewRegistry(eventBus) // TODO: Move outside

// 	top := NewEnvSelector(eventBus, registry)
// 	middle := NewResourceTree(eventBus, registry)

// 	return container.NewBorder(top, nil, nil, nil, middle)
// }

// func NewContentArena(eventBus *mvc.EventBus) fyne.CanvasObject {
// 	return container.NewAppTabs(
// 		container.NewTabItem("hello", container.NewVSplit(
// 			widget.NewLabel("request"),
// 			widget.NewLabel("response"),
// 		)),
// 		container.NewTabItem("bye", container.NewVSplit(
// 			widget.NewLabel("env"),
// 			widget.NewLabel("tmp"),
// 		)),
// 	)
// }
