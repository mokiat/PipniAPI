package widget

import "github.com/mokiat/lacking/ui/state"

func (c *codeAreaComponent) createChangeInsertSegment(segment []rune) state.Change {
	segmentLen := len(segment)
	if segmentLen == 0 {
		return nil
	}
	forward := []state.Action{
		c.createActionInsertSegment(c.cursorRow, c.cursorColumn, segment),
		c.createActionMoveCursor(c.cursorRow, c.cursorColumn+segmentLen),
		c.createActionMoveSelector(c.cursorRow, c.cursorColumn+segmentLen),
	}
	reverse := []state.Action{
		c.createActionMoveSelector(c.selectorRow, c.selectorColumn),
		c.createActionMoveCursor(c.cursorRow, c.cursorColumn),
		c.createActionDeleteSegment(c.cursorRow, c.cursorColumn, c.cursorColumn+segmentLen),
	}
	return c.createChange(forward, reverse)
}

func (c *codeAreaComponent) createChangeInsertLines(lines [][]rune) state.Change {
	lineCount := len(lines)
	if lineCount == 0 {
		return nil
	}
	if lineCount == 1 {
		return c.createChangeInsertSegment(lines[0])
	}
	newCursorRow := c.cursorRow + lineCount - 1
	newCursorColumn := len(lines[lineCount-1])
	forward := []state.Action{
		c.createActionInsertSegment(c.cursorRow, c.cursorColumn, lines[0]),
		c.createActionInsertLines(c.cursorRow+1, lines[1:]),
		c.createActionMoveCursor(newCursorRow, newCursorColumn),
		c.createActionMoveSelector(newCursorRow, newCursorColumn),
	}
	reverse := []state.Action{
		c.createActionMoveSelector(c.selectorRow, c.selectorColumn),
		c.createActionMoveCursor(c.cursorRow, c.cursorColumn),
		c.createActionDeleteLines(c.cursorRow+1, c.cursorRow+lineCount),
		c.createActionDeleteSegment(c.cursorRow, c.cursorColumn, c.cursorColumn+len(lines[0])),
	}
	return c.createChange(forward, reverse)
}

func (c *codeAreaComponent) createChange(forward, reverse []state.Action) state.Change {
	return state.AccumActionChange(forward, reverse, codeAreaChangeAccumulationDuration)
}

func (c *codeAreaComponent) applyChange(change state.Change) {
	c.history.Do(change)
}
