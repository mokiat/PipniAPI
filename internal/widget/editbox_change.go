package widget

import (
	"github.com/mokiat/lacking/ui/state"
	"golang.org/x/exp/slices"
)

func (c *editboxComponent) createChangeInsertSegment(segment []rune) state.Change {
	segmentLen := len(segment)
	if segmentLen == 0 {
		return nil
	}
	forward := []state.Action{
		c.createActionInsertSegment(c.cursorColumn, segment),
		c.createActionMoveCursor(c.cursorColumn + segmentLen),
		c.createActionMoveSelector(c.cursorColumn + segmentLen),
	}
	reverse := []state.Action{
		c.createActionMoveSelector(c.selectorColumn),
		c.createActionMoveCursor(c.cursorColumn),
		c.createActionDeleteSegment(c.cursorColumn, c.cursorColumn+segmentLen),
	}
	return c.createChange(forward, reverse)
}

func (c *editboxComponent) createChangeDeleteSelection() state.Change {
	fromColumn, toColumn := c.selectedColumns()
	if fromColumn >= toColumn {
		return nil
	}
	selectedSegment := slices.Clone(c.line[fromColumn:toColumn])
	forward := []state.Action{
		c.createActionMoveSelector(fromColumn),
		c.createActionMoveCursor(fromColumn),
		c.createActionDeleteSegment(fromColumn, toColumn),
	}
	reverse := []state.Action{
		c.createActionInsertSegment(fromColumn, selectedSegment),
		c.createActionMoveCursor(c.cursorColumn),
		c.createActionMoveSelector(c.selectorColumn),
	}
	return c.createChange(forward, reverse)
}

func (c *editboxComponent) createChangeReplaceSelection(segment []rune) state.Change {
	fromColumn, toColumn := c.selectedColumns()
	if fromColumn >= toColumn {
		return nil
	}
	selectedSegment := slices.Clone(c.line[fromColumn:toColumn])
	forward := []state.Action{
		c.createActionDeleteSegment(fromColumn, toColumn),
		c.createActionInsertSegment(fromColumn, segment),
		c.createActionMoveCursor(fromColumn + len(segment)),
		c.createActionMoveSelector(fromColumn + len(segment)),
	}
	reverse := []state.Action{
		c.createActionMoveCursor(c.cursorColumn),
		c.createActionMoveSelector(c.selectorColumn),
		c.createActionDeleteSegment(fromColumn, fromColumn+len(segment)),
		c.createActionInsertSegment(fromColumn, selectedSegment),
	}
	return c.createChange(forward, reverse)
}

func (c *editboxComponent) createChangeDeleteCharacterLeft() state.Change {
	if c.cursorColumn == 0 {
		return nil
	}
	deletedSegment := slices.Clone(c.line[c.cursorColumn-1 : c.cursorColumn])
	forward := []state.Action{
		c.createActionMoveCursor(c.cursorColumn - 1),
		c.createActionMoveSelector(c.cursorColumn - 1),
		c.createActionDeleteSegment(c.cursorColumn-1, c.cursorColumn),
	}
	reverse := []state.Action{
		c.createActionInsertSegment(c.cursorColumn-1, deletedSegment),
		c.createActionMoveSelector(c.selectorColumn),
		c.createActionMoveCursor(c.cursorColumn),
	}
	return c.createChange(forward, reverse)
}

func (c *editboxComponent) createChangeDeleteCharacterRight() state.Change {
	if c.cursorColumn >= len(c.line) {
		return nil
	}
	deletedSegment := slices.Clone(c.line[c.cursorColumn : c.cursorColumn+1])
	forward := []state.Action{
		c.createActionMoveCursor(c.cursorColumn),
		c.createActionMoveSelector(c.cursorColumn),
		c.createActionDeleteSegment(c.cursorColumn, c.cursorColumn+1),
	}
	reverse := []state.Action{
		c.createActionInsertSegment(c.cursorColumn, deletedSegment),
		c.createActionMoveSelector(c.selectorColumn),
		c.createActionMoveCursor(c.cursorColumn),
	}
	return c.createChange(forward, reverse)
}

func (c *editboxComponent) createChange(forward, reverse []state.Action) state.Change {
	return state.AccumActionChange(forward, reverse, editboxChangeAccumulationDuration)
}

func (c *editboxComponent) applyChange(change state.Change) {
	c.history.Do(change)
}
