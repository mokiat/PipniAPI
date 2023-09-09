package workspace

type Editor interface {
	ID() string
	Name() string

	CanSave() bool
	Save() error

	CanUndo() bool
	Undo()
	CanRedo() bool
	Redo()
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

type NoHistoryEditor struct{}

func (NoHistoryEditor) CanUndo() bool {
	return false
}

func (NoHistoryEditor) Undo() {}

func (NoHistoryEditor) CanRedo() bool {
	return false
}

func (NoHistoryEditor) Redo() {}
