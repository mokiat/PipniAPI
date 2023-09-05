package endpoint

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = co.Define(&editorComponent{})

type EditorData struct {
	EditorModel *endpoint.Editor
}

type editorComponent struct {
	co.BaseComponent
}

func (c *editorComponent) Render() co.Instance {
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
					Items: []std.DropdownItem{
						{
							Key:   http.MethodGet,
							Label: http.MethodGet,
						},
						{
							Key:   http.MethodPost,
							Label: http.MethodPost,
						},
					},
					SelectedKey: http.MethodGet,
				})
			}))

			co.WithChild("uri", co.New(std.Editbox, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
				})
				co.WithData(std.EditboxData{
					Text: "https://example.com",
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
