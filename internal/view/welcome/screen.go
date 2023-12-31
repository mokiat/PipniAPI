package welcome

import (
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"
)

const welcomeMessage = `
Welcome to Pipni API



Here is how to get started:

1. Create an Endpoint resource via the add (+) button in the lower left corner

2. Open the Endpoint resource by selecting it from the list to the left

3. Configure the request data using the Editor that will open here

4. Click 'Go' to execute the request

5. Continue from there
`

var Screen = co.Define(&screenComponent{})

type screenComponent struct {
	co.BaseComponent
}

func (c *screenComponent) Render() co.Instance {
	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			Padding:         ui.UniformSpacing(5),
			BackgroundColor: opt.V(std.SurfaceColor),
			Layout:          layout.Fill(),
		})

		co.WithChild("info", co.New(std.Label, func() {
			co.WithData(std.LabelData{
				Font:      co.OpenFont(c.Scope(), "ui:///roboto-regular.ttf"),
				FontSize:  opt.V(float32(24.0)),
				FontColor: opt.V(std.OnSurfaceColor),
				Text:      welcomeMessage,
			})
		}))
	})
}
