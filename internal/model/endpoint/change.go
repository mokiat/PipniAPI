package endpoint

import (
	"time"

	"github.com/mokiat/PipniAPI/internal/state"
)

func newChangeURIChange(editor *Editor, from, to string) state.Change {
	return &changeURIChange{
		editor:  editor,
		from:    from,
		to:      to,
		created: time.Now(),
	}
}

var _ state.Change = (*changeURIChange)(nil)

var _ state.ExtendableChange = (*changeURIChange)(nil)

type changeURIChange struct {
	editor  *Editor
	from    string
	to      string
	created time.Time
}

func (c *changeURIChange) Apply() {
	c.editor.setURI(c.to)
}

func (c *changeURIChange) Revert() {
	c.editor.setURI(c.from)
}

func (c *changeURIChange) Extend(other state.Change) bool {
	otherURIChange, ok := other.(*changeURIChange)
	if !ok {
		return false
	}
	if otherURIChange.created.Sub(c.created) > 500*time.Millisecond {
		return false
	}
	c.to = otherURIChange.to
	c.created = otherURIChange.created
	return true
}
