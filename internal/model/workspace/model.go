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

	editors    []Editor
	selectedID string
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
		m.SetSelectedID(existing.ID())
		return
	}

	m.editors = append(m.editors, editor)
	m.eventBus.Notify(EditorAddedEvent{
		Model:  m,
		Editor: editor,
	})

	m.SetSelectedID(editor.ID())
}

func (m *Model) RemoveEditor(editor Editor) {
	index := slices.Index(m.editors, editor)
	if index < 0 {
		return
	}

	if m.selectedID == editor.ID() {
		if index > 0 {
			m.SetSelectedID(m.editors[index-1].ID())
		} else if len(m.editors) > 1 {
			m.SetSelectedID(m.editors[1].ID())
		} else {
			m.SetSelectedID("")
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

func (m *Model) SelectedID() string {
	return m.selectedID
}

func (m *Model) SetSelectedID(id string) {
	if id != m.selectedID {
		m.selectedID = id
		m.eventBus.Notify(EditorSelectedEvent{
			Model:  m,
			Editor: m.FindEditor(id),
		})
	}
}

func (m *Model) SelectedEditor() Editor {
	return m.FindEditor(m.selectedID)
}

func (m Model) IsDirty() bool {
	return slices.ContainsFunc(m.editors, func(editor Editor) bool {
		return editor.CanSave()
	})
}
