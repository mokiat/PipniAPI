package endpoint

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var ResponseHeaders = mvc.EventListener(co.Define(&responseHeadersComponent{}))

type ResponseHeadersData struct {
	EditorModel *endpoint.Editor
}

type responseHeadersComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor
}

func (c *responseHeadersComponent) OnUpsert() {
	data := co.GetData[ResponseHeadersData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *responseHeadersComponent) Render() co.Instance {
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

					co.WithChild("name", co.New(std.EditBox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(0),
							Width: opt.V(200),
						})
						co.WithData(std.EditBoxData{
							ReadOnly: true,
							Text:     name,
						})
					}))

					co.WithChild("value", co.New(std.EditBox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(205),
							Right: opt.V(0),
						})
						co.WithData(std.EditBoxData{
							ReadOnly: true,
							Text:     value,
						})
					}))
				}))
			})
		}))
	})
}

func (c *responseHeadersComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case endpoint.ResponseHeadersChangedEvent:
		if event.Editor == c.mdlEditor {
			c.Invalidate()
		}
	}
}

func (c *responseHeadersComponent) eachHeader(cb func(int, string, string)) {
	for i, kv := range c.mdlEditor.ResponseHeaders() {
		cb(i, kv.Key, kv.Value)
	}
}
