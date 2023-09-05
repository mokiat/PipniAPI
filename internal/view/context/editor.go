package context

import (
	"github.com/mokiat/PipniAPI/internal/model/context"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = co.Define(&editorComponent{})

type EditorData struct {
	ContextModel *context.Model
	EditorModel  *context.Editor
}

type editorComponent struct {
	co.BaseComponent
}

func (c *editorComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BackgroundColor: opt.V(ui.Green()),
			BorderSize:      ui.UniformSpacing(1),
			BorderColor:     opt.V(std.OutlineColor),
			Padding:         ui.UniformSpacing(2),
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Top: 5,
				},
			}),
		})
	})
}
