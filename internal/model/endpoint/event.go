package endpoint

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
