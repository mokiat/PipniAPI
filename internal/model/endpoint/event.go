package endpoint

import "net/http"

type MethodChangedEvent struct {
	Editor *Editor
	Method string
}

type URIChangedEvent struct {
	Editor *Editor
	URI    string
}

type RequestBodyChangedEvent struct {
	Editor *Editor
	Body   string
}

type ResponseBodyChangedEvent struct {
	Editor *Editor
	Body   string
}

type ResponseHeadersChangedEvent struct {
	Editor  *Editor
	Headers http.Header
}

type RequestTabChangedEvent struct {
	Editor *Editor
	Tab    EditorTab
}

type ResponseTabChangedEvent struct {
	Editor *Editor
	Tab    EditorTab
}
