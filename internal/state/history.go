package state

import (
	"github.com/mokiat/PipniAPI/internal/pds"
	"github.com/mokiat/gog/ds"
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	historyCapacity = 25
)

func NewHistory(eventBus *mvc.EventBus) *History {
	return &History{
		eventBus: eventBus,

		undoStack: pds.NewClampStack[Change](historyCapacity),
		redoStack: ds.NewStack[Change](historyCapacity),
	}
}

type History struct {
	eventBus *mvc.EventBus

	undoStack *pds.ClampStack[Change]
	redoStack *ds.Stack[Change]
}

func (h *History) LastChange() Change {
	if h.undoStack.IsEmpty() {
		return nil
	}
	return h.undoStack.Peek()
}

func (h *History) Do(change Change) {
	h.undoStack.Push(change)
	h.redoStack.Clear()
	change.Apply()
	h.notifyChanged()
}

func (h *History) CanUndo() bool {
	return !h.undoStack.IsEmpty()
}

func (h *History) Undo() {
	change := h.undoStack.Pop()
	h.redoStack.Push(change)
	change.Apply()
	h.notifyChanged()
}

func (h *History) CanRedo() bool {
	return !h.redoStack.IsEmpty()
}

func (h *History) Redo() {
	change := h.redoStack.Pop()
	h.undoStack.Push(change)
	change.Apply()
	h.notifyChanged()
}

func (h *History) notifyChanged() {
	h.eventBus.Notify(HistoryChangedEvent{
		History: h,
	})
}

type HistoryChangedEvent struct {
	History *History
}
