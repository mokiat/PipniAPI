package state

import (
	"github.com/mokiat/PipniAPI/internal/pds"
	"github.com/mokiat/gog/ds"
)

const (
	historyCapacity = 25
)

func NewHistory() *History {
	return &History{
		undoStack: pds.NewClampStack[Change](historyCapacity),
		redoStack: ds.NewStack[Change](historyCapacity),
	}
}

type History struct {
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
	if extChange, ok := h.extendableChange(); !ok || !extChange.Extend(change) {
		h.undoStack.Push(change)
	}
	h.redoStack.Clear()
	change.Apply()
}

func (h *History) CanUndo() bool {
	return !h.undoStack.IsEmpty()
}

func (h *History) Undo() {
	change := h.undoStack.Pop()
	h.redoStack.Push(change)
	change.Revert()
}

func (h *History) CanRedo() bool {
	return !h.redoStack.IsEmpty()
}

func (h *History) Redo() {
	change := h.redoStack.Pop()
	h.undoStack.Push(change)
	change.Apply()
}

func (h *History) extendableChange() (ExtendableChange, bool) {
	if h.undoStack.IsEmpty() {
		return nil, false
	}
	change := h.undoStack.Peek()
	extendable, ok := change.(ExtendableChange)
	if !ok {
		return nil, false
	}
	return extendable, true
}
