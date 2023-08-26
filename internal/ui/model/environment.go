package model

func NewEnvironmentEditor() Editor {
	return &EnvironmentEditor{}
}

type EnvironmentEditor struct {
}

func (e *EnvironmentEditor) ID() string {
	return "environment"
}

func (e *EnvironmentEditor) Title() string {
	return "Environment"
}
