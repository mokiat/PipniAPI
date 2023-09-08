package endpoint

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, endpoint *registry.Endpoint) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		endpoint: endpoint,

		// TODO: Initialize the following from mdlEndpoint once.
		method: http.MethodGet,
		uri:    "https://api.publicapis.org/entries",
	}
}

type Editor struct {
	workspace.NoSaveEditor
	workspace.NoHistoryEditor

	eventBus *mvc.EventBus
	endpoint *registry.Endpoint

	method       string
	uri          string
	requestBody  string
	responseBody string
}

func (e *Editor) ID() string {
	return e.endpoint.ID()
}

func (e *Editor) Name() string {
	return e.endpoint.Name()
}

func (e *Editor) Method() string {
	return e.method
}

func (e *Editor) SetMethod(method string) {
	if method != e.method {
		e.method = method
		e.eventBus.Notify(MethodChangedEvent{
			Editor: e,
			Method: method,
		})
	}
}

func (e *Editor) URI() string {
	return e.uri
}

func (e *Editor) SetURI(uri string) {
	if uri != e.uri {
		e.uri = uri
		e.eventBus.Notify(URIChangedEvent{
			Editor: e,
			URI:    uri,
		})
	}
}

func (e *Editor) RequestBody() string {
	return e.requestBody
}

func (e *Editor) SetRequestBody(body string) {
	if body != e.requestBody {
		e.requestBody = body
		e.eventBus.Notify(RequestBodyChangedEvent{
			Editor: e,
			Body:   body,
		})
	}
}

func (e *Editor) ResponseBody() string {
	return e.responseBody
}

func (e *Editor) SetResponseBody(body string) {
	if body != e.responseBody {
		e.responseBody = body
		e.eventBus.Notify(ResponseBodyChangedEvent{
			Editor: e,
			Body:   body,
		})
	}
}