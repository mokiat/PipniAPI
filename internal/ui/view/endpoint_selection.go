package view

import (
	"github.com/mokiat/PipniAPI/internal/ui/model"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var EndpointSelection = mvc.EventListener(co.Define(&endpointSelectionComponent{}))

type EndpointSelectionData struct {
	RegistryModel *model.Registry
}

type endpointSelectionComponent struct {
	co.BaseComponent

	mdlRegistry *model.Registry
}

func (c *endpointSelectionComponent) OnUpsert() {
	data := co.GetData[EndpointSelectionData](c.Properties())
	c.mdlRegistry = data.RegistryModel
}

func (c *endpointSelectionComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.SymmetricSpacing(0, 1),
			BorderColor: opt.V(std.OutlineColor),
			Layout:      layout.Fill(),
			Padding:     ui.SymmetricSpacing(0, 5),
		})

		co.WithChild("list", co.New(std.List, func() {

			for _, endpoint := range c.mdlRegistry.Root().Endpoints() {
				endpoint := endpoint
				co.WithChild(endpoint.ID(), co.New(EndpointItem, func() {
					co.WithData(EndpointItemData{
						Selected: c.mdlRegistry.SelectedID() == endpoint.ID(),
						Icon:     co.OpenImage(c.Scope(), "images/ping.png"),
						Text:     endpoint.Name(),
					})
					co.WithCallbackData(EndpointItemCallbackData{
						OnClick: func() {
							c.onEndpointSelected(endpoint)
						},
					})
				}))
			}

			for _, workflow := range c.mdlRegistry.Root().Workflows() {
				workflow := workflow
				co.WithChild(workflow.ID(), co.New(EndpointItem, func() {
					co.WithData(EndpointItemData{
						Selected: c.mdlRegistry.SelectedID() == workflow.ID(),
						Icon:     co.OpenImage(c.Scope(), "images/workflow.png"),
						Text:     workflow.Name(),
					})
					co.WithCallbackData(EndpointItemCallbackData{
						OnClick: func() {
							c.onWorkflowSelected(workflow)
						},
					})
				}))
			}
		}))
	})
}

func (c *endpointSelectionComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case model.RegistrySelectionChangedEvent:
		c.Invalidate()
	}
}

func (c *endpointSelectionComponent) onEndpointSelected(endpoint *model.Endpoint) {
	c.mdlRegistry.SetSelectedID(endpoint.ID())
}

func (c *endpointSelectionComponent) onWorkflowSelected(workflow *model.Workflow) {
	c.mdlRegistry.SetSelectedID(workflow.ID())
}

var EndpointItem = co.Define(&endpointItemComponent{})

type EndpointItemData struct {
	Selected bool
	Icon     *ui.Image
	Text     string
}

type EndpointItemCallbackData struct {
	OnClick std.OnActionFunc
}

type endpointItemComponent struct {
	co.BaseComponent

	selected bool
	icon     *ui.Image
	text     string

	onClick std.OnActionFunc
}

func (c *endpointItemComponent) OnUpsert() {
	data := co.GetData[EndpointItemData](c.Properties())
	c.selected = data.Selected
	c.icon = data.Icon
	c.text = data.Text

	callbackData := co.GetCallbackData[EndpointItemCallbackData](c.Properties())
	c.onClick = callbackData.OnClick
}

func (c *endpointItemComponent) Render() co.Instance {
	return co.New(std.ListItem, func() {
		co.WithLayoutData(layout.Data{
			GrowHorizontally: true,
		})
		co.WithData(std.ListItemData{
			Selected: c.selected,
		})
		co.WithCallbackData(std.ListItemCallbackData{
			OnSelected: c.onClick,
		})

		co.WithChild("holder", co.New(std.Element, func() {
			co.WithData(std.ElementData{
				Layout: layout.Frame(),
			})

			if c.icon != nil {
				co.WithChild("icon", co.New(std.Picture, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentLeft,
						VerticalAlignment:   layout.VerticalAlignmentCenter,
						Width:               opt.V(24),
						Height:              opt.V(24),
					})
					co.WithData(std.PictureData{
						Image:      c.icon,
						ImageColor: opt.V(std.OnSurfaceColor),
						Mode:       std.ImageModeFit,
					})
				}))
			}

			if c.text != "" {
				co.WithChild("label", co.New(std.Label, func() {
					co.WithLayoutData(layout.Data{
						HorizontalAlignment: layout.HorizontalAlignmentCenter,
					})
					co.WithData(std.LabelData{
						Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
						FontSize:  opt.V(float32(20.0)),
						FontColor: opt.V(std.OnSurfaceColor),
						Text:      c.text,
					})
				}))
			}
		}))
	})
}