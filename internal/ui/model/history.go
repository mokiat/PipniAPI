package model

import (
	"github.com/mokiat/PipniAPI/internal/pds"
	"github.com/mokiat/gog/ds"
	"github.com/mokiat/lacking/ui/mvc"
)

const (
	historyCapacity = 25
)

type Change interface {
	Apply()
	Revert()
}

func FuncChange(apply, revert func()) Change {
	return &funcChange{
		apply:  apply,
		revert: revert,
	}
}

type funcChange struct {
	apply  func()
	revert func()
}

func (ch *funcChange) Apply() {
	ch.apply()
}

func (ch *funcChange) Revert() {
	ch.revert()
}

func CombinedChange(changes ...Change) Change {
	return &combinedChange{
		changes: changes,
	}
}

type combinedChange struct {
	changes []Change
}

func (c *combinedChange) Apply() {
	for i := 0; i < len(c.changes); i++ {
		c.changes[i].Apply()
	}
}

func (c *combinedChange) Revert() {
	for i := len(c.changes) - 1; i >= 0; i-- {
		c.changes[i].Revert()
	}
}

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
