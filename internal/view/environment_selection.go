package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"github.com/mokiat/gog"
)

func (w *Window) newEnvironmentSelection() fyne.CanvasObject {
	envSelectControl := widget.NewSelect(nil, nil)

	updateOptions := func(envs []*model.Environment) {
		values := gog.Map(w.mdlRegistry.Environments(), func(env *model.Environment) string {
			return env.Name()
		})
		envSelectControl.Options = values
		envSelectControl.Refresh()
	}
	updateOptions(w.mdlRegistry.Environments())

	updateSelected := func(activeEnv *model.Environment) {
		if activeEnv != nil {
			envSelectControl.SetSelected(activeEnv.Name())
		} else {
			envSelectControl.SetSelected("")
		}
	}
	updateSelected(w.mdlRegistry.ActiveEnvironment())

	envSelectControl.OnChanged = func(name string) {
		env, ok := gog.FindFunc(w.mdlRegistry.Environments(), func(env *model.Environment) bool {
			return env.Name() == name
		})
		if ok {
			w.mdlRegistry.SetActiveEnvironment(env)
		}
	}

	settingsButton := widget.NewButton("Settings", nil)
	settingsButton.OnTapped = func() {
		// FIXME: Prevent duplicates
		w.mdlWorkspace.AddEditor(model.NewEnvironmentEditor())
	}

	w.eventBus.Subscribe(func(event mvc.Event) {
		switch event := event.(type) {
		case model.ActiveEnvironmentChangedEvent:
			updateSelected(event.ActiveEnvironment)
		}
	})

	return container.NewBorder(nil, nil, nil, settingsButton, envSelectControl)
}
