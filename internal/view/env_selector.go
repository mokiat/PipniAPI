package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"github.com/mokiat/gog"
)

func NewEnvSelector(
	eventBus *mvc.EventBus,
	mdlRegistry *model.Registry,
	mdlWorkspace *model.Workspace,
) fyne.CanvasObject {

	envSelectControl := widget.NewSelect(nil, nil)

	updateOptions := func(envs []*model.Environment) {
		values := gog.Map(mdlRegistry.Environments(), func(env *model.Environment) string {
			return env.Name()
		})
		envSelectControl.Options = values
		envSelectControl.Refresh()
	}
	updateOptions(mdlRegistry.Environments())

	updateSelected := func(activeEnv *model.Environment) {
		if activeEnv != nil {
			envSelectControl.SetSelected(activeEnv.Name())
		} else {
			envSelectControl.SetSelected("")
		}
	}
	updateSelected(mdlRegistry.ActiveEnvironment())

	envSelectControl.OnChanged = func(name string) {
		env, ok := gog.FindFunc(mdlRegistry.Environments(), func(env *model.Environment) bool {
			return env.Name() == name
		})
		if ok {
			mdlRegistry.SetActiveEnvironment(env)
		}
	}

	settingsButton := widget.NewButton("Settings", nil)
	settingsButton.OnTapped = func() {
		// FIXME: Prevent duplicates
		mdlWorkspace.AddEditor(model.NewEnvironmentEditor())
	}

	eventBus.Subscribe(func(event mvc.Event) {
		switch event := event.(type) {
		case model.ActiveEnvironmentChangedEvent:
			updateSelected(event.ActiveEnvironment)
		}
	})

	return container.NewBorder(nil, nil, nil, settingsButton, envSelectControl)
}
