package endpoint

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var supportedMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

var Editor = mvc.EventListener(co.Define(&editorComponent{}))

type EditorData struct {
	EditorModel *endpoint.Editor
}

type editorComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *editorComponent) Render() co.Instance {
	methodItems := gog.Map(supportedMethods, func(method string) std.DropdownItem {
		return std.DropdownItem{
			Key:   method,
			Label: method,
		}
	})

	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.UniformSpacing(1),
			BorderColor: opt.V(std.OutlineColor),
			Padding:     ui.UniformSpacing(2),
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Top: 5,
				},
			}),
		})

		co.WithChild("uri-settings", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(std.ElementData{
				Layout: layout.Frame(layout.FrameSettings{
					ContentSpacing: ui.SymmetricSpacing(5, 0),
				}),
			})

			co.WithChild("method", co.New(std.Dropdown, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(100),
				})
				co.WithData(std.DropdownData{
					Items:       methodItems,
					SelectedKey: c.mdlEditor.Method(),
				})
				co.WithCallbackData(std.DropdownCallbackData{
					OnItemSelected: func(key any) {
						c.changeMethod(key.(string))
					},
				})
			}))

			co.WithChild("uri", co.New(std.Editbox, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
				})
				co.WithData(std.EditboxData{
					Text: c.mdlEditor.URI(),
				})
			}))

			co.WithChild("go", co.New(std.Button, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(100),
				})
				co.WithData(std.ButtonData{
					Text: "GO", // TODO: Fix button to align text in the middle
				})
			}))
		}))

		co.WithChild("payload-settings", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ContainerData{
				Layout: layout.Frame(),
			})

			co.WithChild("request", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
					Width:               opt.V(2000), // FIXME: Use different layout
				})
				co.WithData(std.ContainerData{
					Padding:         ui.UniformSpacing(5),
					BorderColor:     opt.V(std.OutlineColor),
					BorderSize:      ui.UniformSpacing(1),
					BackgroundColor: opt.V(ui.Gray()),
				})
			}))

			co.WithChild("response", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ContainerData{
					Padding:         ui.UniformSpacing(5),
					BorderColor:     opt.V(std.OutlineColor),
					BorderSize:      ui.UniformSpacing(1),
					BackgroundColor: opt.V(ui.Gray()),
				})
			}))
		}))
	})
}

func (c *editorComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case endpoint.MethodChangedEvent:
		c.Invalidate()
	}
}

func (c *editorComponent) changeMethod(method string) {
	c.mdlEditor.ChangeMethod(method)
}
