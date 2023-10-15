package endpoint

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"text/template"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, reg *registry.Model, endpoint *registry.Endpoint) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		reg:      reg,
		endpoint: endpoint,

		method:          endpoint.Method(),
		uri:             endpoint.URI(),
		requestBody:     endpoint.Body(),
		requestHeaders:  endpoint.Headers(),
		responseBody:    "",
		responseHeaders: nil,

		requestTab:  EditorTabBody,
		responseTab: EditorTabBody,
	}
}

type Editor struct {
	eventBus *mvc.EventBus
	reg      *registry.Model
	endpoint *registry.Endpoint

	method string
	uri    string

	requestTab     EditorTab
	requestBody    string
	requestHeaders []gog.KV[string, string]

	responseTab     EditorTab
	responseBody    string
	responseHeaders []gog.KV[string, string]
}

func (e *Editor) ID() string {
	return e.endpoint.ID()
}

func (e *Editor) Name() string {
	return e.endpoint.Name()
}

func (e *Editor) CanSave() bool {
	if e.method != e.endpoint.Method() {
		return true
	}
	if e.uri != e.endpoint.URI() {
		return true
	}
	if !slices.Equal(e.requestHeaders, e.endpoint.Headers()) {
		return true
	}
	if e.requestBody != e.endpoint.Body() {
		return true
	}
	return false
}

func (e *Editor) Save() error {
	e.endpoint.SetMethod(e.method)
	e.endpoint.SetURI(e.uri)
	e.endpoint.SetHeaders(e.requestHeaders)
	e.endpoint.SetBody(e.requestBody)
	if err := e.reg.Save(); err != nil {
		return fmt.Errorf("error saving registry: %w", err)
	}
	e.notifyModified()
	return nil
}

func (e *Editor) Method() string {
	return e.method
}

func (e *Editor) SetMethod(method string) {
	if method != e.method {
		e.method = method
		e.eventBus.Notify(MethodChangedEvent{
			Editor: e,
		})
		e.notifyModified()
	}
}

func (e *Editor) URI() string {
	return e.uri
}

func (e *Editor) HTTPURI() string {
	return e.evaluate(e.uri)
}

func (e *Editor) SetURI(newURI string) {
	if newURI != e.uri {
		e.uri = newURI
		e.eventBus.Notify(URIChangedEvent{
			Editor: e,
		})
		e.notifyModified()
	}
}

func (e *Editor) RequestHeaders() []gog.KV[string, string] {
	return slices.Clone(e.requestHeaders)
}

func (e *Editor) HTTPRequestHeaders() http.Header {
	result := make(http.Header)
	for _, header := range e.requestHeaders {
		result.Add(e.evaluate(header.Key), e.evaluate(header.Value))
	}
	return result
}

func (e *Editor) AddRequestHeader() {
	e.requestHeaders = append(e.requestHeaders, gog.KV[string, string]{
		Key:   "",
		Value: "",
	})
	e.eventBus.Notify(RequestHeadersChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) SetRequestHeaderName(index int, name string) {
	e.requestHeaders[index].Key = name
	e.eventBus.Notify(RequestHeadersChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) SetRequestHeaderValue(index int, value string) {
	e.requestHeaders[index].Value = value
	e.eventBus.Notify(RequestHeadersChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) DeleteRequestHeader(index int) {
	e.requestHeaders = slices.Delete(e.requestHeaders, index, index+1)
	e.eventBus.Notify(RequestHeadersChangedEvent{
		Editor: e,
	})
	e.notifyModified()
}

func (e *Editor) RequestBody() string {
	return e.requestBody
}

func (e *Editor) HTTPRequestBody() string {
	return e.evaluate(e.requestBody)
}

func (e *Editor) SetRequestBody(body string) {
	if body != e.requestBody {
		e.requestBody = body
		e.eventBus.Notify(RequestBodyChangedEvent{
			Editor: e,
		})
		e.notifyModified()
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
		})
	}
}

func (e *Editor) ResponseHeaders() []gog.KV[string, string] {
	return e.responseHeaders
}

func (e *Editor) SetHTTPResponseHeaders(headers http.Header) {
	e.responseHeaders = e.responseHeaders[:0]
	for name, values := range headers {
		for _, value := range values {
			e.responseHeaders = append(e.responseHeaders, gog.KV[string, string]{
				Key:   name,
				Value: value,
			})
		}
	}
	e.eventBus.Notify(ResponseHeadersChangedEvent{
		Editor: e,
	})
}

func (e *Editor) RequestTab() EditorTab {
	return e.requestTab
}

func (e *Editor) SetRequestTab(tab EditorTab) {
	if tab != e.requestTab {
		e.requestTab = tab
		e.eventBus.Notify(RequestTabChangedEvent{
			Editor: e,
		})
	}
}

func (e *Editor) ResponseTab() EditorTab {
	return e.responseTab
}

func (e *Editor) SetResponseTab(tab EditorTab) {
	if tab != e.responseTab {
		e.responseTab = tab
		e.eventBus.Notify(ResponseTabChangedEvent{
			Editor: e,
		})
	}
}

func (e *Editor) evaluate(text string) string {
	// FIXME: This evaluation is slow. It tries to find the context
	// for each evaluation and is too tightly coupled. This should be
	// managed differently (through an API elsewhere).

	activeContext := e.reg.ActiveContext()
	if activeContext == nil {
		return text
	}

	tmpl, err := template.New("tmp").Parse(text)
	if err != nil {
		log.Warn("Cannot evaluate text expression: %v", err)
		return text
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, activeContext.DataProperties()); err != nil {
		log.Warn("Cannot evaluate text expression: %v", err)
		return text
	}
	return output.String()
}

func (e *Editor) notifyModified() {
	e.eventBus.Notify(workspace.EditorModifiedEvent{
		Editor: e,
	})
}

type EditorTab string

const (
	EditorTabBody    EditorTab = "body"
	EditorTabHeaders EditorTab = "headers"
	EditorTabStats   EditorTab = "stats"
)
