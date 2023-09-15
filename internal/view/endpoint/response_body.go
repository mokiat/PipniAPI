package endpoint

import (
	"github.com/mokiat/PipniAPI/internal/view/widget"
	co "github.com/mokiat/lacking/ui/component"
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
	return co.New(widget.CodeArea, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(widget.CodeAreaData{
			ReadOnly: true,
			Code:     c.text,
		})
	})
}
