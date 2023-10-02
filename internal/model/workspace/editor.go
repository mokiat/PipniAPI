package workspace

type Editor interface {
	ID() string
	Name() string

	CanSave() bool
	Save() error
}

type EditorModifiedEvent struct {
	Editor Editor
}

type NoSaveEditor struct{}

func (NoSaveEditor) CanSave() bool {
	return false
}

func (NoSaveEditor) Save() error {
	return nil
}
