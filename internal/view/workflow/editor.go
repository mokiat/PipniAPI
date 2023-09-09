package workflow

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

var Editor = co.Define(&editorComponent{})

type EditorData struct{}

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
			Layout:      layout.Anchor(),
		})

		co.WithChild("under-dev", co.New(std.Label, func() {
			co.WithLayoutData(layout.Data{
				HorizontalCenter: opt.V(0),
				VerticalCenter:   opt.V(0),
			})
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-bold.ttf"),
				FontSize:  opt.V(float32(32)),
				FontColor: opt.V(std.OnSurfaceColor),
				Text:      "Under Development",
			})
		}))
	})
}
