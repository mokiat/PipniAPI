package endpoint

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var RequestHeaders = mvc.EventListener(co.Define(&requestHeadersComponent{}))

type RequestHeadersData struct {
	EditorModel *endpoint.Editor
}

type requestHeadersComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor
}

func (c *requestHeadersComponent) OnUpsert() {
	data := co.GetData[RequestHeadersData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *requestHeadersComponent) Render() co.Instance {
	return co.New(std.ScrollPane, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ScrollPaneData{
			DisableHorizontal: true,
		})

		co.WithChild("table", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				GrowHorizontally: true,
			})
			co.WithData(std.ElementData{
				Padding: ui.UniformSpacing(5),
				Layout: layout.Vertical(layout.VerticalSettings{
					ContentSpacing: 5,
				}),
			})

			c.eachHeader(func(index int, name, value string) {
				co.WithChild(fmt.Sprintf("row-%d", index), co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Layout: layout.Anchor(),
					})

					co.WithChild("name", co.New(widget.EditBox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(0),
							Width: opt.V(200),
						})
						co.WithData(widget.EditBoxData{
							Text: name,
						})
						co.WithCallbackData(widget.EditBoxCallbackData{
							OnChange: func(newName string) {
								c.changeHeaderName(index, newName)
							},
						})
					}))

					co.WithChild("value", co.New(widget.EditBox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(205),
							Right: opt.V(40),
						})
						co.WithData(widget.EditBoxData{
							Text: value,
						})
						co.WithCallbackData(widget.EditBoxCallbackData{
							OnChange: func(newValue string) {
								c.changeHeaderValue(index, newValue)
							},
						})
					}))

					co.WithChild("delete", co.New(std.Button, func() {
						co.WithLayoutData(layout.Data{
							Right: opt.V(0),
							Width: opt.V(36),
						})
						co.WithData(std.ButtonData{
							Icon: co.OpenImage(c.Scope(), "images/delete.png"),
							Text: value,
						})
						co.WithCallbackData(std.ButtonCallbackData{
							OnClick: func() {
								c.deleteHeader(index)
							},
						})
					}))
				}))
			})

			co.WithChild("add-header-button", co.New(std.Button, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
				})
				co.WithData(std.ButtonData{
					Icon: co.OpenImage(c.Scope(), "images/add.png"),
					Text: "Add Header",
				})
				co.WithCallbackData(std.ButtonCallbackData{
					OnClick: func() {
						c.addHeader()
					},
				})
			}))
		}))
	})
}

func (c *requestHeadersComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case endpoint.RequestHeadersChangedEvent:
		if event.Editor == c.mdlEditor {
			c.Invalidate()
		}
	}
}

func (c *requestHeadersComponent) eachHeader(cb func(int, string, string)) {
	for i, kv := range c.mdlEditor.RequestHeaders() {
		cb(i, kv.Key, kv.Value)
	}
}

func (c *requestHeadersComponent) addHeader() {
	c.mdlEditor.AddRequestHeader()
}

func (c *requestHeadersComponent) changeHeaderName(index int, name string) {
	c.mdlEditor.SetRequestHeaderName(index, name)
}

func (c *requestHeadersComponent) changeHeaderValue(index int, value string) {
	c.mdlEditor.SetRequestHeaderValue(index, value)
}

func (c *requestHeadersComponent) deleteHeader(index int) {
	c.mdlEditor.DeleteRequestHeader(index)
}
