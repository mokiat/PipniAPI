package endpoint

import (
	"cmp"
	"net/http"
	"slices"

	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var ResponseHeaders = co.Define(&responseHeadersComponent{})

type ResponseHeadersData struct {
	Headers http.Header
}

type responseHeadersComponent struct {
	co.BaseComponent

	headers http.Header
}

func (c *responseHeadersComponent) OnUpsert() {
	data := co.GetData[ResponseHeadersData](c.Properties())
	c.headers = data.Headers
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

			c.eachHeader(func(name, value string) {
				co.WithChild("row", co.New(std.Element, func() {
					co.WithLayoutData(layout.Data{
						GrowHorizontally: true,
					})
					co.WithData(std.ElementData{
						Layout: layout.Anchor(),
					})

					co.WithChild("name", co.New(std.Editbox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(0),
							Width: opt.V(200),
						})
						co.WithData(std.EditboxData{
							Text: name,
							// TODO: Readonly
						})
					}))

					co.WithChild("value", co.New(std.Editbox, func() {
						co.WithLayoutData(layout.Data{
							Left:  opt.V(205),
							Right: opt.V(0),
						})
						co.WithData(std.EditboxData{
							Text: value,
							// TODO: Readonly
						})
					}))
				}))
			})
		}))
	})
}

func (c *responseHeadersComponent) eachHeader(cb func(string, string)) {
	headerEntries := gog.Entries(c.headers)
	slices.SortFunc(headerEntries, func(a, b gog.KV[string, []string]) int {
		return cmp.Compare(a.Key, b.Key)
	})
	for _, entry := range headerEntries {
		name := entry.Key
		for _, value := range entry.Value {
			cb(name, value)
		}
	}
}
