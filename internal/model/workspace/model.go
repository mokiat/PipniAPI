package workspace

import (
	"github.com/mokiat/lacking/ui/mvc"
	"golang.org/x/exp/slices"
)

func NewModel(eventBus *mvc.EventBus) *Model {
	return &Model{
		eventBus: eventBus,
	}
}

type Model struct {
	eventBus *mvc.EventBus

	editors        []Editor
	selectedEditor Editor
}

func (m *Model) Editors() []Editor {
	return m.editors
}

func (m *Model) EachEditor(cb func(editor Editor)) {
	for _, editor := range m.editors {
		cb(editor)
	}
}

func (m *Model) FindEditor(id string) Editor {
	for _, editor := range m.editors {
		if editor.ID() == id {
			return editor
		}
	}
	return nil
}

func (m *Model) AppendEditor(editor Editor) {
	if existing := m.FindEditor(editor.ID()); existing != nil {
		m.SelectEditor(existing)
		return
	}

	m.editors = append(m.editors, editor)
	m.eventBus.Notify(EditorAddedEvent{
		Model:  m,
		Editor: editor,
	})

	m.SelectEditor(editor)
}

func (m *Model) RemoveEditor(editor Editor) {
	index := slices.Index(m.editors, editor)
	if index < 0 {
		return
	}

	if m.selectedEditor == editor {
		if index > 0 {
			m.SelectEditor(m.editors[index-1])
		} else if len(m.editors) > 1 {
			m.SelectEditor(m.editors[1])
		} else {
			m.SelectEditor(nil)
		}
	}

	m.editors = slices.DeleteFunc(m.editors, func(candidate Editor) bool {
		return candidate == editor
	})
	m.eventBus.Notify(EditorRemovedEvent{
		Model:  m,
		Editor: editor,
	})
}

func (m *Model) SelectedEditor() Editor {
	return m.selectedEditor
}

func (m *Model) SelectEditor(editor Editor) {
	if editor != m.selectedEditor {
		m.selectedEditor = editor
		m.eventBus.Notify(EditorSelectedEvent{
			Model:  m,
			Editor: editor,
		})
	}
}
