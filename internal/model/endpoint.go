package model

import (
	"net/http"

	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewEndpointEditor(eventBus *mvc.EventBus, id string) Editor {
	return &EndpointEditor{
		eventBus: eventBus,

		id:     id,
		name:   "Example",
		method: http.MethodGet,
		uri:    "https://api.publicapis.org/entries",
	}
}

type EndpointEditor struct {
	eventBus *mvc.EventBus

	id     string
	name   string
	method string

	uri string
}

func (e *EndpointEditor) ID() string {
	return e.id
}

func (e *EndpointEditor) Title() string {
	return e.name
}

func (e *EndpointEditor) Name() string {
	return e.name
}

func (e *EndpointEditor) SetName(name string) {
	if name != e.name {
		e.name = name
		e.eventBus.Notify(EndpointNameChangedEvent{
			Editor: e,
			Name:   name,
		})
	}
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

type EndpointNameChangedEvent struct {
	Editor *EndpointEditor
	Name   string
}

type EndpointMethodChangedEvent struct {
	Editor *EndpointEditor
	Method string
}

type EndpointURIChangedEvent struct {
	Editor *EndpointEditor
	URI    string
}
