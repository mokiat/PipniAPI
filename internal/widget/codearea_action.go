package widget

import (
	"github.com/mokiat/gog"
	"github.com/mokiat/lacking/ui/state"
	"golang.org/x/exp/slices"
)

func (c *codeAreaComponent) createActionMoveCursor(row, column int) state.Action {
	return func() {
		c.cursorRow = row
		c.cursorColumn = column
	}
}

func (c *codeAreaComponent) createActionMoveSelector(row, column int) state.Action {
	return func() {
		c.selectorRow = row
		c.selectorColumn = column
	}
}

func (c *codeAreaComponent) createActionInsertSegment(row, column int, segment []rune) state.Action {
	return func() {
		if len(segment) > 0 {
			line := c.lines[row]
			prefix := line[:column]
			suffix := line[column:]
			c.lines[row] = gog.Concat(prefix, segment, suffix)
		}
	}
}

func (c *codeAreaComponent) createActionDeleteSegment(row, fromColumn, toColumn int) func() {
	return func() {
		if fromColumn < toColumn {
			c.lines[row] = slices.Delete(c.lines[row], fromColumn, toColumn)
		}
	}
}

func (c *codeAreaComponent) createActionInsertLines(row int, lines [][]rune) func() {
	return func() {
		if len(lines) > 0 {
			c.lines = slices.Insert(c.lines, row, slices.Clone(lines)...)
		}
	}
}

func (c *codeAreaComponent) createActionDeleteLines(fromRow, toRow int) func() {
	return func() {
		if fromRow < toRow {
			c.lines = slices.Delete(c.lines, fromRow, toRow)
		}
	}
}
