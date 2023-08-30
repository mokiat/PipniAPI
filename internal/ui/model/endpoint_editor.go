package model

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/registrymodel"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEndpointEditor(eventBus *mvc.EventBus, endpoint *registrymodel.Endpoint) Editor {
	return &EndpointEditor{
		eventBus: eventBus,
		endpoint: endpoint,
		history:  NewHistory(eventBus),

		// TODO: Initialize the following from mdlEndpoint once.
		method: http.MethodGet,
		uri:    "https://api.publicapis.org/entries",
	}
}

type EndpointEditor struct {
	eventBus *mvc.EventBus
	endpoint *registrymodel.Endpoint
	history  *History

	method       string
	uri          string
	requestBody  string
	responseBody string
}

func (e *EndpointEditor) History() *History {
	return e.history
}

func (e *EndpointEditor) ID() string {
	return e.endpoint.ID()
}

func (e *EndpointEditor) Title() string {
	return e.endpoint.Name()
}

func (e *EndpointEditor) Method() string {
	return e.method
}

func (e *EndpointEditor) SetMethod(method string) {
	if method != e.method {
		e.method = method
		e.eventBus.Notify(EndpointMethodChangedEvent{
			Editor: e,
			Method: method,
		})
	}
}

func (e *EndpointEditor) URI() string {
	return e.uri
}

func (e *EndpointEditor) SetURI(uri string) {
	if uri != e.uri {
		e.uri = uri
		e.eventBus.Notify(EndpointURIChangedEvent{
			Editor: e,
			URI:    uri,
		})
	}
}

func (e *EndpointEditor) RequestBody() string {
	return e.requestBody
}

func (e *EndpointEditor) SetRequestBody(body string) {
	e.requestBody = body
	// TODO: Notify
}

func (e *EndpointEditor) ResponseBody() string {
	return e.responseBody
}

func (e *EndpointEditor) SetResponseBody(body string) {
	if body != e.responseBody {
		e.responseBody = body
		e.eventBus.Notify(EndpointResponseBodyChangedEvent{
			Editor: e,
			Body:   body,
		})
	}
}

type EndpointMethodChangedEvent struct {
	Editor *EndpointEditor
	Method string
}

type EndpointURIChangedEvent struct {
	Editor *EndpointEditor
	URI    string
}

type EndpointResponseBodyChangedEvent struct {
	Editor *EndpointEditor
	Body   string
}
