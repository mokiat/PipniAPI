package widget

import (
	"github.com/mokiat/gog"
	"golang.org/x/exp/slices"
)

func (c *editboxComponent) createActionMoveCursor(column int) func() {
	return func() {
		c.cursorColumn = column
	}
}

func (c *editboxComponent) createActionMoveSelector(column int) func() {
	return func() {
		c.selectorColumn = column
	}
}

func (c *editboxComponent) createActionInsertSegment(column int, segment []rune) func() {
	return func() {
		prefix := c.line[:column]
		suffix := c.line[column:]
		c.line = gog.Concat(prefix, segment, suffix)
	}
}

func (c *editboxComponent) createActionDeleteSegment(fromColumn, toColumn int) func() {
	return func() {
		c.line = slices.Delete(c.line, fromColumn, toColumn)
	}
}
