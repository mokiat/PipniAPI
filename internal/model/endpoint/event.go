package endpoint

type MethodChangedEvent struct {
	Editor *Editor
}

type URIChangedEvent struct {
	Editor *Editor
}

type RequestTabChangedEvent struct {
	Editor *Editor
}

type RequestBodyChangedEvent struct {
	Editor *Editor
}

type RequestHeadersChangedEvent struct {
	Editor *Editor
}

type ResponseTabChangedEvent struct {
	Editor *Editor
}

type ResponseBodyChangedEvent struct {
	Editor *Editor
}

type ResponseHeadersChangedEvent struct {
	Editor *Editor
}
