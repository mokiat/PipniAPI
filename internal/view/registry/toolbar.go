package registry

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/debug/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Toolbar = mvc.EventListener(co.Define(&toolbarComponent{}))

type ToolbarData struct {
	RegistryModel *registry.Model
}

type toolbarComponent struct {
	co.BaseComponent

	mdlRegistry *registry.Model
}

func (c *toolbarComponent) OnUpsert() {
	data := co.GetData[ToolbarData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *toolbarComponent) Render() co.Instance {
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
				OnClick: func() {
					c.openEditResourceModal(resource)
				},
			})
		}))

		co.WithChild("clone", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/duplicate.png"),
				Enabled: opt.V(canClone),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: func() {
					c.cloneResource(resource)
				},
			})
		}))

		co.WithChild("delete", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/delete.png"),
				Enabled: opt.V(canDelete),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: func() {
					c.openDeleteResourceModal(resource)
				},
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
				OnClick: func() {
					c.moveResourceUp(resource)
				},
			})
		}))

		co.WithChild("move-down", co.New(std.Button, func() {
			co.WithData(std.ButtonData{
				Icon:    co.OpenImage(c.Scope(), "images/move-down.png"),
				Enabled: opt.V(canMoveDown),
			})
			co.WithCallbackData(std.ButtonCallbackData{
				OnClick: func() {
					c.moveResourceDown(resource)
				},
			})
		}))
	})
}

func (c *toolbarComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case registry.RegistrySelectionChangedEvent:
		c.Invalidate()
	case registry.RegistryStructureChangedEvent:
		c.Invalidate()
	}
}

func (c *toolbarComponent) openAddResourceModal() {
	co.OpenOverlay(c.Scope(), co.New(Modal, func() {
		co.WithData(ModalData{
			Name:          "",
			Kind:          registry.ResourceKindEndpoint,
			CanChangeKind: true,
		})
		co.WithCallbackData(ModalCallbackData{
			OnApply: c.addResource,
		})
	}))
}

func (c *toolbarComponent) openEditResourceModal(resource registry.Resource) {
	co.OpenOverlay(c.Scope(), co.New(Modal, func() {
		co.WithData(ModalData{
			Name:          resource.Name(),
			Kind:          resource.Kind(),
			CanChangeKind: false,
		})
		co.WithCallbackData(ModalCallbackData{
			OnApply: func(name string, _ registry.ResourceKind) {
				c.renameResource(resource, name)
			},
		})
	}))
}

func (c *toolbarComponent) openDeleteResourceModal(resource registry.Resource) {
	co.OpenOverlay(c.Scope(), co.New(widget.ConfirmationModal, func() {
		co.WithData(widget.ConfirmationModalData{
			Icon: co.OpenImage(c.Scope(), "images/warning.png"),
			Text: fmt.Sprintf("Are you sure you want to delete the following resource:\n\n\n%q\n\n\nThis cannot be undone!", resource.Name()),
		})
		co.WithCallbackData(widget.ConfirmationModalCallbackData{
			OnApply: func() {
				c.deleteResource(resource)
			},
		})
	}))
}

func (c *toolbarComponent) addResource(name string, kind registry.ResourceKind) {
	c.mdlRegistry.CreateResource(c.mdlRegistry.Root(), name, kind)
	c.saveChanges()
}

func (c *toolbarComponent) renameResource(resource registry.Resource, name string) {
	c.mdlRegistry.RenameResource(resource, name)
	c.saveChanges()
}

func (c *toolbarComponent) cloneResource(resource registry.Resource) {
	c.mdlRegistry.CloneResource(resource)
	c.saveChanges()
}

func (c *toolbarComponent) deleteResource(resource registry.Resource) {
	c.mdlRegistry.DeleteResource(resource)
	c.saveChanges()
}

func (c *toolbarComponent) moveResourceUp(resource registry.Resource) {
	c.mdlRegistry.MoveUp(resource)
	c.saveChanges()
}

func (c *toolbarComponent) moveResourceDown(resource registry.Resource) {
	c.mdlRegistry.MoveDown(resource)
	c.saveChanges()
}

func (c *toolbarComponent) saveChanges() {
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
