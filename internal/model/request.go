package model

func NewRequestEditor() Editor {
	return &RequestEditor{}
}

type RequestEditor struct {
	_forceUnique bool
}

func (e *RequestEditor) Title() string {
	return "Create User" // TODO
}
