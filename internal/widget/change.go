package widget

import (
	"time"

	"github.com/mokiat/lacking/ui/state"
)

const (
	textTypeAccumulationDuration = 500 * time.Millisecond
)

func emptyTextTypeChange() *textTypeChange {
	return &textTypeChange{
		when: time.Now(),
	}
}

type textTypeChange struct {
	when    time.Time
	forward []func()
	reverse []func()
}

func (c *textTypeChange) Apply() {
	for _, action := range c.forward {
		action()
	}
}

func (c *textTypeChange) Revert() {
	for _, action := range c.reverse {
		action()
	}
}

func (c *textTypeChange) Extend(other state.Change) bool {
	otherChange, ok := other.(*textTypeChange)
	if !ok {
		return false
	}
	if otherChange.when.Sub(c.when) > textTypeAccumulationDuration {
		return false
	}
	c.forward = append(c.forward, otherChange.forward...)
	c.reverse = append(otherChange.reverse, c.reverse...)
	c.when = otherChange.when
	return true
}
