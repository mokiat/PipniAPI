package context

import (
	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Selector = mvc.EventListener(co.Define(&selectorComponent{}))

type SelectorData struct {
	WorkspaceModel *workspace.Model
	ContextModel   *context.Model
}

type selectorComponent struct {
	co.BaseComponent

	eventBus     *mvc.EventBus
	mdlWorkspace *workspace.Model
	mdlContext   *context.Model
}

func (c *selectorComponent) OnUpsert() {
	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	data := co.GetData[SelectorData](c.Properties())
	c.mdlWorkspace = data.WorkspaceModel
	c.mdlContext = data.ContextModel
}

func (c *selectorComponent) Render() co.Instance {
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

func (c *selectorComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case context.EnvironmentSelectedEvent:
		c.Invalidate()
	}
}

func (c *selectorComponent) onDropdownItemSelected(key any) {
	c.mdlContext.SetSelectedID(key.(string))
	c.saveChanges()
}

func (c *selectorComponent) onSettingsClicked() {
	if editor := c.mdlWorkspace.FindEditor(context.EditorID); editor != nil {
		c.mdlWorkspace.SetSelectedID(editor.ID())
	} else {
		c.mdlWorkspace.AppendEditor(context.NewEditor(c.eventBus, c.mdlContext))
	}
}

func (c *selectorComponent) saveChanges() {
	if err := c.mdlContext.Save(); err != nil {
		log.Error("Error saving context: %v", err)
		co.OpenOverlay(c.Scope(), co.New(widget.NotificationModal, func() {
			co.WithData(widget.NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: "The program encountered an error.\n\nChanges could not be saved.\n\nCheck logs for more information.",
			})
		}))
	}
}
