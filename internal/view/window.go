package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

type Window struct {
	app fyne.App
	win fyne.Window

	eventBus     *mvc.EventBus
	mdlRegistry  *model.Registry
	mdlWorkspace *model.Workspace
}

func NewWindow(app fyne.App, win fyne.Window) *Window {
	eventBus := mvc.NewEventBus()

	mdlRegistry := model.NewRegistry(eventBus)

	mdlWorkspace := model.NewWorkspace(eventBus)
	mdlWorkspace.AddEditor(model.NewEnvironmentEditor())
	mdlWorkspace.AddEditor(model.NewEndpointEditor(eventBus, "<guid-here>"))

	return &Window{
		app: app,
		win: win,

		eventBus:     eventBus,
		mdlRegistry:  mdlRegistry,
		mdlWorkspace: mdlWorkspace,
	}
}

func (w *Window) RenderMainMenu() *fyne.MainMenu {
	return w.newMainMenu()
}

func (w *Window) RenderContent() fyne.CanvasObject {
	return container.NewHSplit(
		w.newNavigationPanel(),
		w.newWorkspace(),
	)
}
