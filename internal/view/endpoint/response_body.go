package endpoint

import (
	"github.com/mokiat/gog/opt"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
)

var ResponseBody = co.Define(&responseBodyComponent{})

type ResponseBodyData struct {
	Text string
}

type responseBodyComponent struct {
	co.BaseComponent

	text string
}

func (c *responseBodyComponent) OnUpsert() {
	data := co.GetData[ResponseBodyData](c.Properties())
	c.text = data.Text
}

func (c *responseBodyComponent) Render() co.Instance {
	// TODO: TextArea
	return co.New(std.Label, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.LabelData{
			Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
			FontSize:  opt.V(float32(24.0)),
			FontColor: opt.V(std.OnSurfaceColor),
			Text:      c.text,
		})
	})
}
