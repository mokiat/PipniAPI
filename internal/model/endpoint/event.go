package endpoint

import "github.com/mokiat/gog"

type MethodChangedEvent struct {
	Editor *Editor
	Method string
}

type URIChangedEvent struct {
	Editor *Editor
	URI    string
}

type RequestTabChangedEvent struct {
	Editor *Editor
	Tab    EditorTab
}

type RequestBodyChangedEvent struct {
	Editor *Editor
	Body   string
}

type RequestHeadersChangedEvent struct {
	Editor  *Editor
	Headers []gog.KV[string, string]
}

type ResponseTabChangedEvent struct {
	Editor *Editor
	Tab    EditorTab
}

type ResponseBodyChangedEvent struct {
	Editor *Editor
	Body   string
}

type ResponseHeadersChangedEvent struct {
	Editor  *Editor
	Headers []gog.KV[string, string]
}
