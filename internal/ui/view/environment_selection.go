package view

import (
	"github.com/mokiat/PipniAPI/internal/ui/model"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var EnvironmentSelection = mvc.EventListener(co.Define(&environmentSelectionComponent{}))

type EnvironmentSelectionData struct {
	ContextModel *model.Context
}

type environmentSelectionComponent struct {
	co.BaseComponent

	mdlContext *model.Context
}

func (c *environmentSelectionComponent) OnUpsert() {
	data := co.GetData[EnvironmentSelectionData](c.Properties())
	c.mdlContext = data.ContextModel
}

func (c *environmentSelectionComponent) Render() co.Instance {
	dropdownItems := gog.Map(c.mdlContext.Environments(), func(env *model.Environment) std.DropdownItem {
		return std.DropdownItem{
			Key:   env.ID(),
			Label: env.Name(),
		}
	})
	selectedItem := c.mdlContext.SelectedID()

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
				SelectedKey: selectedItem,
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
		}))
	})
}

func (c *environmentSelectionComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.ContextSelectionChangedEvent:
		c.Invalidate()
	}
}

func (c *environmentSelectionComponent) onDropdownItemSelected(key any) {
	c.mdlContext.SetSelectedID(key.(string))
}
