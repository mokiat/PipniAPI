package context

import (
	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/debug/log"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Selector = mvc.EventListener(co.Define(&selectorComponent{}))

type SelectorData struct {
	RegistryModel *registry.Model
}

type selectorComponent struct {
	co.BaseComponent

	eventBus    *mvc.EventBus
	mdlRegistry *registry.Model
}

func (c *selectorComponent) OnUpsert() {
	c.eventBus = co.TypedValue[*mvc.EventBus](c.Scope())

	data := co.GetData[SelectorData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *selectorComponent) Render() co.Instance {
	contexts := c.listContexts()
	dropdownItems := gog.Map(contexts, func(env *registry.Context) std.DropdownItem {
		return std.DropdownItem{
			Key:   env.ID(),
			Label: env.Name(),
		}
	})

	var selectedKey string
	if activeContext := c.mdlRegistry.ActiveContext(); activeContext != nil {
		selectedKey = activeContext.ID()
	}

	return co.New(std.Dropdown, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.DropdownData{
			Items:       dropdownItems,
			SelectedKey: selectedKey,
		})
		co.WithCallbackData(std.DropdownCallbackData{
			OnItemSelected: c.onDropdownItemSelected,
		})
	})
}

func (c *selectorComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case registry.RegistryActiveContextChangedEvent:
		c.Invalidate()
	}
}

func (c *selectorComponent) onDropdownItemSelected(key any) {
	c.mdlRegistry.SetActiveContextID(key.(string))
	c.saveChanges()
}

func (c *selectorComponent) listContexts() []*registry.Context {
	resources := c.mdlRegistry.AllResources()
	contextResources := gog.Select(resources, func(resource registry.Resource) bool {
		_, ok := resource.(*registry.Context)
		return ok
	})
	return gog.Map(contextResources, func(resource registry.Resource) *registry.Context {
		return resource.(*registry.Context)
	})
}

func (c *selectorComponent) saveChanges() {
	if err := c.mdlRegistry.Save(); err != nil {
		log.Error("Error saving registry: %v", err)
		co.OpenOverlay(c.Scope(), co.New(widget.NotificationModal, func() {
			co.WithData(widget.NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: "The program encountered an error.\n\nChanges could not be saved.\n\nCheck logs for more information.",
			})
		}))
	}
}
