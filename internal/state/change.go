package state

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
