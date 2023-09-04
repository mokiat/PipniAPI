package view

import (
	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var EnvironmentSelection = mvc.EventListener(co.Define(&environmentSelectionComponent{}))

type EnvironmentSelectionData struct {
	WorkspaceModel *workspace.Model
	ContextModel   *context.Model
}

type environmentSelectionComponent struct {
	co.BaseComponent

	eventBus     *mvc.EventBus
	mdlWorkspace *workspace.Model
	mdlContext   *context.Model
}

func (c *environmentSelectionComponent) OnUpsert() {
	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	data := co.GetData[EnvironmentSelectionData](c.Properties())
	c.mdlWorkspace = data.WorkspaceModel
	c.mdlContext = data.ContextModel
}

func (c *environmentSelectionComponent) Render() co.Instance {
	dropdownItems := gog.Map(c.mdlContext.Environments(), func(env *context.Environment) std.DropdownItem {
		return std.DropdownItem{
			Key:   env.ID(),
			Label: env.Name(),
		}
	})
	var selectedKey string
	if selectedEnv := c.mdlContext.SelectedEnvironment(); selectedEnv != nil {
		selectedKey = selectedEnv.ID()
	}

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Right: 5,
				},
			}),
		})

		co.WithChild("dropdown", co.New(std.Dropdown, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.DropdownData{
				Items:       dropdownItems,
				SelectedKey: selectedKey,
			})
			co.WithCallbackData(std.DropdownCallbackData{
				OnItemSelected: c.onDropdownItemSelected,
			})
		}))

		co.WithChild("settings", co.New(std.Button, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentRight,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/settings.png"),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.onSettingsClicked,
			})
		}))
	})
}

func (c *environmentSelectionComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case context.EnvironmentSelectedEvent:
		c.Invalidate()
	}
}

func (c *environmentSelectionComponent) onDropdownItemSelected(key any) {
	environment := c.mdlContext.FindEnvironment(key.(string))
	c.mdlContext.SelectEnvironment(environment)
}

func (c *environmentSelectionComponent) onSettingsClicked() {
	if editor := c.mdlWorkspace.FindEditor(context.ContextEditorID); editor != nil {
		c.mdlWorkspace.SelectEditor(editor)
	} else {
		c.mdlWorkspace.AppendEditor(context.NewEditor(c.eventBus, c.mdlContext))
	}
}
