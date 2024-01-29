package endpoint

import (
	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var ResponseBody = mvc.EventListener(co.Define(&responseBodyComponent{}))

type ResponseBodyData struct {
	EditorModel *endpoint.Editor
}

type responseBodyComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor
}

func (c *responseBodyComponent) OnUpsert() {
	data := co.GetData[ResponseBodyData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *responseBodyComponent) Render() co.Instance {
	return co.New(std.CodeArea, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.CodeAreaData{
			ReadOnly: true,
			Code:     c.mdlEditor.ResponseBody(),
		})
	})
}

func (c *responseBodyComponent) OnEvent(event mvc.Event) {
	switch event := event.(type) {
	case endpoint.ResponseBodyChangedEvent:
		if event.Editor == c.mdlEditor {
			c.Invalidate()
		}
	}
}
