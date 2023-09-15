package endpoint

import (
	"github.com/mokiat/PipniAPI/internal/view/widget"
	co "github.com/mokiat/lacking/ui/component"
)

var RequestBody = co.Define(&requestBodyComponent{})

type RequestBodyData struct {
	Text string
}

type RequestBodyCallbackData struct {
	OnChange func(text string)
}

type requestBodyComponent struct {
	co.BaseComponent

	text string

	onChange func(string)
}

func (c *requestBodyComponent) OnUpsert() {
	data := co.GetData[RequestBodyData](c.Properties())
	c.text = data.Text

	callbackData := co.GetCallbackData[RequestBodyCallbackData](c.Properties())
	c.onChange = callbackData.OnChange
}

func (c *requestBodyComponent) Render() co.Instance {
	return co.New(widget.CodeArea, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(widget.CodeAreaData{
			Code: c.text,
		})
		co.WithCallbackData(widget.CodeAreaCallbackData{
			OnChange: c.onChange,
		})
	})
}
