package context

import (
	"github.com/mokiat/PipniAPI/internal/model/context"
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

	})
}

func (c *editorComponent) OnEvent(event mvc.Event) {

}
