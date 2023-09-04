package view

import (
	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var EndpointManagement = mvc.EventListener(co.Define(&endpointManagementComponent{}))

type EndpointManagementData struct {
	RegistryModel *registrymodel.Registry
}

type endpointManagementComponent struct {
	co.BaseComponent

	mdlRegistry *registrymodel.Registry
}

func (c *endpointManagementComponent) OnUpsert() {
	data := co.GetData[EndpointManagementData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *endpointManagementComponent) Render() co.Instance {
	resource := c.mdlRegistry.FindSelectedResource()
	canEdit := resource != nil
	canClone := resource != nil
	canDelete := resource != nil
	canMoveUp := c.mdlRegistry.CanMoveUp(resource)
	canMoveDown := c.mdlRegistry.CanMoveDown(resource)

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Layout: layout.Horizontal(layout.HorizontalSettings{
				ContentAlignment: layout.VerticalAlignmentCenter,
				ContentSpacing:   2,
			}),
		})

		co.WithChild("add", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon: co.OpenImage(c.Scope(), "images/add.png"),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.openAddResourceModal,
			})
		}))

		co.WithChild("edit", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/edit.png"),
				Enabled: opt.V(canEdit),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.editResource,
			})
		}))

		co.WithChild("clone", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/duplicate.png"),
				Enabled: opt.V(canClone),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.cloneResource,
			})
		}))

		co.WithChild("delete", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/delete.png"),
				Enabled: opt.V(canDelete),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.deleteResource,
			})
		}))

		co.WithChild("separator", co.New(std.Spacing, func() {
			co.WithData(std.SpacingData{
				Size: ui.NewSize(10, 0),
			})
		}))

		co.WithChild("move-up", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/move-up.png"),
				Enabled: opt.V(canMoveUp),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.moveSelectionUp,
			})
		}))

		co.WithChild("move-down", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/move-down.png"),
				Enabled: opt.V(canMoveDown),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: c.moveSelectionDown,
			})
		}))
	})
}

func (c *endpointManagementComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case registrymodel.RegistrySelectionChangedEvent:
		c.Invalidate()
	case registrymodel.RegistryStructureChangedEvent:
		c.Invalidate()
	}
}

func (c *endpointManagementComponent) openAddResourceModal() {
	co.OpenOverlay(c.Scope(), co.New(ResourceModal, func() {
		co.WithData(ResourceModalData{
			Name:          "",
			Kind:          registrymodel.ResourceKindEndpoint,
			CanChangeKind: true,
		})
		co.WithCallbackData(ResourceModalCallbackData{
			OnApply: c.addResource,
		})
	}))
}

func (c *endpointManagementComponent) addResource(name string, kind registrymodel.ResourceKind) {
	c.mdlRegistry.CreateResource(c.mdlRegistry.Root(), name, kind)
	if err := c.mdlRegistry.Save(); err != nil {
		panic(err) // TODO: Display error message
	}
}

func (c *endpointManagementComponent) editResource() {
	resource := c.mdlRegistry.FindSelectedResource()

	co.OpenOverlay(c.Scope(), co.New(ResourceModal, func() {
		co.WithData(ResourceModalData{
			Name:          resource.Name(),
			Kind:          resource.Kind(),
			CanChangeKind: false,
		})
		co.WithCallbackData(ResourceModalCallbackData{
			OnApply: func(name string, _ registrymodel.ResourceKind) {
				c.renameResource(resource, name)
			},
		})
	}))
}

func (c *endpointManagementComponent) renameResource(resource registrymodel.Resource, name string) {
	c.mdlRegistry.RenameResource(resource, name)
	if err := c.mdlRegistry.Save(); err != nil {
		panic(err) // TODO: Display error message
	}
}

func (c *endpointManagementComponent) cloneResource() {

}

func (c *endpointManagementComponent) deleteResource() {

}

func (c *endpointManagementComponent) moveSelectionUp() {
	c.mdlRegistry.MoveUp(c.mdlRegistry.FindSelectedResource())
	if err := c.mdlRegistry.Save(); err != nil {
		panic(err) // TODO: Display error message
	}
}

func (c *endpointManagementComponent) moveSelectionDown() {
	c.mdlRegistry.MoveDown(c.mdlRegistry.FindSelectedResource())
	if err := c.mdlRegistry.Save(); err != nil {
		panic(err) // TODO: Display error message
	}
}
