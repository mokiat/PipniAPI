package model

type Editor interface {
	Title() string
}

// func NewEditor(eventBus *mvc.EventBus) *Editor {
// 	return &Editor{
// 		eventBus: eventBus,
// 	}
// }

// type Editor struct {
// 	eventBus *mvc.EventBus
// }

// func (e *Editor) IsDirty() bool {
// 	return false
// }

// func (e *Editor) Save() {

// }

// type EditorDirtyChangedEvent struct {
// 	Editor *Editor
// }
