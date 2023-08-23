package model

func NewEnvironmentEditor() Editor {
	return &EnvironmentEditor{}
}

type EnvironmentEditor struct {
	_forceUnique bool
}

func (e *EnvironmentEditor) Title() string {
	return "Environment"
}
