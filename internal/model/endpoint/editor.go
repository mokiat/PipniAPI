package endpoint

import (
	"fmt"
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/registry"
	"github.com/mokiat/PipniAPI/internal/model/workspace"
	"github.com/mokiat/PipniAPI/internal/state"
	"github.com/mokiat/lacking/ui/mvc"
)

func NewEditor(eventBus *mvc.EventBus, reg *registry.Model, endpoint *registry.Endpoint) workspace.Editor {
	return &Editor{
		eventBus: eventBus,
		reg:      reg,
		endpoint: endpoint,

		history:     state.NewHistory(),
		savedChange: nil,

		method:          endpoint.Method(),
		uri:             endpoint.URI(),
		requestHeaders:  endpoint.Headers(),
		requestBody:     endpoint.Body(),
		responseHeaders: make(http.Header),
		responseBody:    "",

		requestTab:  EditorTabBody,
		responseTab: EditorTabBody,
	}
}

type Editor struct {
	eventBus *mvc.EventBus
	reg      *registry.Model
	endpoint *registry.Endpoint

	history     *state.History
	savedChange state.Change

	method          string
	uri             string
	requestHeaders  http.Header
	requestBody     string
	responseHeaders http.Header
	responseBody    string

	requestTab  EditorTab
	responseTab EditorTab
}

func (e *Editor) ID() string {
	return e.endpoint.ID()
}

func (e *Editor) Name() string {
	return e.endpoint.Name()
}

func (e *Editor) CanSave() bool {
	return e.savedChange != e.history.LastChange()
}

func (e *Editor) Save() error {
	e.endpoint.SetMethod(e.method)
	e.endpoint.SetURI(e.uri)
	e.endpoint.SetHeaders(e.requestHeaders)
	e.endpoint.SetBody(e.requestBody)
	if err := e.reg.Save(); err != nil {
		return fmt.Errorf("error saving registry: %w", err)
	}
	e.savedChange = e.history.LastChange()
	e.notifyModified()
	return nil
}

func (e *Editor) CanUndo() bool {
	return e.history.CanUndo()
}

func (e *Editor) Undo() {
	e.history.Undo()
	e.notifyModified()
}

func (e *Editor) CanRedo() bool {
	return e.history.CanRedo()
}

func (e *Editor) Redo() {
	e.history.Redo()
	e.notifyModified()
}

func (e *Editor) Method() string {
	return e.method
}

func (e *Editor) setMethod(method string) {
	if method != e.method {
		e.method = method
		e.eventBus.Notify(MethodChangedEvent{
			Editor: e,
			Method: method,
		})
	}
}

func (e *Editor) ChangeMethod(newMethod string) {
	oldMethod := e.method
	if newMethod == oldMethod {
		return
	}
	e.do(func() {
		e.setMethod(newMethod)
	}, func() {
		e.setMethod(oldMethod)
	})
}

func (e *Editor) URI() string {
	return e.uri
}

func (e *Editor) ChangeURI(newURI string) {
	oldURI := e.uri
	if newURI == oldURI {
		return
	}
	e.do(func() {
		e.setURI(newURI)
	}, func() {
		e.setURI(oldURI)
	})
}

func (e *Editor) setURI(uri string) {
	if uri != e.uri {
		e.uri = uri
		e.eventBus.Notify(URIChangedEvent{
			Editor: e,
			URI:    uri,
		})
	}
}

func (e *Editor) RequestHeaders() http.Header {
	return e.requestHeaders.Clone()
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

func (e *Editor) ResponseHeaders() http.Header {
	return e.responseHeaders
}

func (e *Editor) SetResponseHeaders(headers http.Header) {
	e.responseHeaders = headers
	e.eventBus.Notify(ResponseHeadersChangedEvent{
		Editor:  e,
		Headers: headers,
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
			Tab:    tab,
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
			Tab:    tab,
		})
	}
}

func (e *Editor) do(apply, revert func()) {
	e.history.Do(state.FuncChange(apply, revert))
	e.notifyModified()
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
