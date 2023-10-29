package context

import (
	"fmt"

	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = mvc.EventListener(co.Define(&editorComponent{}))

type EditorData struct {
	EditorModel *context.Editor
}

type editorComponent struct {
	co.BaseComponent

	mdlEditor *context.Editor
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *editorComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.UniformSpacing(1),
			BorderColor: opt.V(std.OutlineColor),
			Padding:     ui.UniformSpacing(2),
			Layout:      layout.Fill(),
		})

		co.WithChild("scroll-pane", co.New(std.ScrollPane, func() {
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

				c.eachProperty(func(index int, name, value string) {
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
									c.changePropertyName(index, newName)
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
									c.changePropertyValue(index, newValue)
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
									c.deleteProperty(index)
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
						Text: "Add Property",
					})
					co.WithCallbackData(std.ButtonCallbackData{
						OnClick: func() {
							c.addProperty()
						},
					})
				}))
			}))
		}))
	})
}

func (c *editorComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case context.PropertiesChangedEvent:
		if event.Editor == c.mdlEditor {
			c.Invalidate()
		}
	}
}

func (c *editorComponent) eachProperty(cb func(int, string, string)) {
	for i, kv := range c.mdlEditor.Properties() {
		cb(i, kv.Key, kv.Value)
	}
}

func (c *editorComponent) addProperty() {
	c.mdlEditor.AddProperty()
}

func (c *editorComponent) changePropertyName(index int, name string) {
	c.mdlEditor.SetPropertyName(index, name)
}

func (c *editorComponent) changePropertyValue(index int, value string) {
	c.mdlEditor.SetPropertyValue(index, value)
}

func (c *editorComponent) deleteProperty(index int) {
	c.mdlEditor.DeleteProperty(index)
}
