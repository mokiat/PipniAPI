package view

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"github.com/mokiat/gog"
)

func NewEnvSelector(eventBus *mvc.EventBus, mdl *model.Registry) fyne.CanvasObject {
	envSelectControl := widget.NewSelect(nil, nil)

	updateOptions := func(envs []*model.Environment) {
		values := gog.Map(mdl.Environments(), func(env *model.Environment) string {
			return env.Name()
		})
		envSelectControl.Options = values
		envSelectControl.Refresh()
	}
	updateOptions(mdl.Environments())

	updateSelected := func(activeEnv *model.Environment) {
		if activeEnv != nil {
			envSelectControl.SetSelected(activeEnv.Name())
		} else {
			envSelectControl.SetSelected("")
		}
	}
	updateSelected(mdl.ActiveEnvironment())

	envSelectControl.OnChanged = func(name string) {
		env, ok := gog.FindFunc(mdl.Environments(), func(env *model.Environment) bool {
			return env.Name() == name
		})
		if ok {
			mdl.SetActiveEnvironment(env)
		}
	}

	settingsButton := widget.NewButton("Settings", nil)
	settingsButton.OnTapped = func() {
		log.Println("SETTINGS")
	}

	eventBus.Subscribe(func(event mvc.Event) {
		switch event := event.(type) {
		case model.ActiveEnvironmentChangedEvent:
			updateSelected(event.ActiveEnvironment)
		}
	})

	return container.NewBorder(nil, nil, nil, settingsButton, envSelectControl)
}
