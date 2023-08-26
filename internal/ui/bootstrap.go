package ui

import (
	"github.com/mokiat/PipniAPI/internal/ui/view"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
)

func BootstrapApplication(window *ui.Window) {
	eventBus := mvc.NewEventBus()

	scope := co.RootScope(window)
	scope = co.TypedValueScope(scope, eventBus)
	co.Initialize(scope, co.New(Bootstrap, nil))
}

var Bootstrap = co.Define(&bootstrapComponent{})

type bootstrapComponent struct {
	co.BaseComponent

	// appModel      *model.Application
	// settingsModel *model.Settings
	// loadingModel  *model.Loading
	// playModel     *model.Play
}

func (c *bootstrapComponent) OnCreate() {
	// eventBus := co.TypedValue[*mvc.EventBus](c.Scope())
	// c.appModel = model.NewApplication(eventBus)
	// c.settingsModel = model.NewSettings(eventBus)
	// c.loadingModel = model.NewLoading()
	// c.playModel = model.NewPlay()
}

func (c *bootstrapComponent) Render() co.Instance {
	return co.New(view.Application, func() {
		// co.WithData(view.ApplicationData{
		// 	// AppModel:      c.appModel,
		// 	// SettingsModel: c.settingsModel,
		// 	// LoadingModel:  c.loadingModel,
		// 	// PlayModel:     c.playModel,
		// })
	})
}
