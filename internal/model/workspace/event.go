package workspace

type EditorAddedEvent struct {
	Model  *Model
	Editor Editor
}

type EditorRemovedEvent struct {
	Model  *Model
	Editor Editor
}

type EditorSelectedEvent struct {
	Model  *Model
	Editor Editor
}
