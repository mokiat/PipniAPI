package widget

import "github.com/mokiat/lacking/ui/state"

func (c *codeAreaComponent) createChange(forward, reverse []state.Action) state.Change {
	return state.AccumActionChange(forward, reverse, codeAreaChangeAccumulationDuration)
}
