package endpoint

import (
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var RequestBody = mvc.EventListener(co.Define(&requestBodyComponent{}))

type RequestBodyData struct {
	EditorModel *endpoint.Editor
}

type requestBodyComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor
}

func (c *requestBodyComponent) OnUpsert() {
	data := co.GetData[RequestBodyData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *requestBodyComponent) Render() co.Instance {
	return co.New(std.CodeArea, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.CodeAreaData{
			Code: c.mdlEditor.RequestBody(),
		})
		co.WithCallbackData(std.CodeAreaCallbackData{
			OnChange: c.changeRequestBody,
		})
	})
}

func (c *requestBodyComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case endpoint.ResponseBodyChangedEvent:
		if event.Editor == c.mdlEditor {
			c.Invalidate()
		}
	}
}

func (c *requestBodyComponent) changeRequestBody(body string) {
	c.mdlEditor.SetRequestBody(body)
}
