package widget

import (
	"time"

	"github.com/mokiat/PipniAPI/internal/shortcuts"
	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/state"
	"github.com/mokiat/lacking/ui/std"
	"golang.org/x/exp/slices"
)

const (
	editBoxHistoryCapacity = 100
)

var EditBox = co.Define(&editBoxComponent{})

type EditBoxData struct {
	ReadOnly bool
	Text     string
}

type EditBoxCallbackData struct {
	OnChange func(string)
	OnSubmit func(string)
}

var _ ui.ElementResizeHandler = (*editBoxComponent)(nil)
var _ ui.ElementRenderHandler = (*editBoxComponent)(nil)
var _ ui.ElementKeyboardHandler = (*editBoxComponent)(nil)
var _ ui.ElementMouseHandler = (*editBoxComponent)(nil)
var _ ui.ElementHistoryHandler = (*editBoxComponent)(nil)
var _ ui.ElementClipboardHandler = (*editBoxComponent)(nil)

type editBoxComponent struct {
	co.BaseComponent

	history *state.History

	font     *ui.Font
	fontSize float32

	cursorColumn   int
	selectorColumn int

	isReadOnly bool
	line       []rune
	onChange   func(string)
	onSubmit   func(string)

	offsetX    int
	maxOffsetX int

	isDragging bool
}

func (c *editBoxComponent) OnCreate() {
	c.history = state.NewHistory(editBoxHistoryCapacity)

	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = 18.0

	c.cursorColumn = 0
}

func (c *editBoxComponent) OnUpsert() {
	data := co.GetData[EditBoxData](c.Properties())
	c.isReadOnly = data.ReadOnly
	if data.Text != string(c.line) {
		c.history.Clear()
	}
	c.line = []rune(data.Text)

	callbackData := co.GetOptionalCallbackData(c.Properties(), EditBoxCallbackData{})
	c.onChange = callbackData.OnChange
	c.onSubmit = callbackData.OnSubmit

	c.cursorColumn = min(c.cursorColumn, len(c.line))
	c.selectorColumn = min(c.selectorColumn, len(c.line))
}

func (c *editBoxComponent) Render() co.Instance {
	contentSize := c.font.TextSize(string(c.line), c.fontSize)
	padding := ui.Spacing{
		Left:   10,
		Right:  10,
		Top:    5,
		Bottom: 5,
	}

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			IdealSize: opt.V(ui.Size{
				Width:  int(contentSize.X) + padding.Horizontal(),
				Height: int(c.fontSize) + padding.Vertical(),
			}),
			Padding: padding,
		})
	})
}

func (c *editBoxComponent) OnResize(element *ui.Element, bounds ui.Bounds) {
	availableTextWidth := bounds.Width - element.Padding().Horizontal()
	textWidth := int(c.font.TextSize(string(c.line), c.fontSize).X)
	c.maxOffsetX = max(textWidth-availableTextWidth, 0)
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editBoxComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	bounds := canvas.DrawBounds(element, false)
	paddedBounds := canvas.DrawBounds(element, true)
	isFocused := element.Window().IsElementFocused(element)

	// Background
	canvas.Reset()
	canvas.RoundRectangle(
		bounds.Position,
		bounds.Size,
		sprec.NewVec4(8, 8, 8, 8),
	)
	canvas.Fill(ui.Fill{
		Color: std.SurfaceColor,
	})

	canvas.Push()
	canvas.SetClipRect(
		paddedBounds.X(),
		paddedBounds.X()+paddedBounds.Width(),
		paddedBounds.Y(),
		paddedBounds.Y()+paddedBounds.Height(),
	)

	textSize := c.font.TextSize(string(c.line), c.fontSize)
	textPosition := sprec.Vec2Diff(paddedBounds.Position, sprec.NewVec2(float32(c.offsetX), 0.0))

	// Draw Selection
	if c.cursorColumn != c.selectorColumn {
		fromColumn := min(c.cursorColumn, c.selectorColumn)
		toColumn := max(c.cursorColumn, c.selectorColumn)
		preSelectionSize := c.font.TextSize(string(c.line[:fromColumn]), c.fontSize)
		selectionSize := c.font.TextSize(string(c.line[fromColumn:toColumn]), c.fontSize)

		selectionPosition := sprec.Vec2Sum(textPosition, sprec.NewVec2(preSelectionSize.X, 0.0))
		canvas.Reset()
		canvas.Rectangle(selectionPosition, selectionSize)
		if isFocused {
			canvas.Fill(ui.Fill{
				Color: std.SecondaryLightColor,
			})
		} else {
			canvas.Fill(ui.Fill{
				Color: ui.MixColors(std.SecondaryLightColor, std.SurfaceColor, 0.75),
			})
		}
	}

	// Draw text
	canvas.Reset()
	canvas.FillText(string(c.line), textPosition, ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: std.OnSurfaceColor,
	})

	// Draw cursor
	if !c.isReadOnly && isFocused {
		preCursorText := c.line[:c.cursorColumn]
		preCursorTextSize := c.font.TextSize(string(preCursorText), c.fontSize)

		// TODO: Take tilt into account and use like stroke instead of rect fill.
		cursorPosition := sprec.Vec2Sum(textPosition, sprec.NewVec2(preCursorTextSize.X, 0.0))
		cursorWidth := float32(1.0)
		canvas.Reset()
		canvas.Rectangle(cursorPosition, sprec.NewVec2(cursorWidth, textSize.Y))
		canvas.Fill(ui.Fill{
			Color: std.PrimaryColor,
		})
	}

	canvas.Pop()

	// Highlight
	canvas.Reset()
	if isFocused {
		canvas.SetStrokeColor(std.SecondaryLightColor)
	} else {
		canvas.SetStrokeColor(std.PrimaryLightColor)
	}
	canvas.SetStrokeSize(2.0)
	canvas.RoundRectangle(
		bounds.Position,
		bounds.Size,
		sprec.NewVec4(8, 8, 8, 8),
	)
	canvas.Stroke()
}

func (c *editBoxComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Action {
	case ui.KeyboardActionDown, ui.KeyboardActionRepeat:
		consumed := c.onKeyboardPressEvent(element, event)
		if consumed {
			element.Invalidate()
		}
		return consumed

	case ui.KeyboardActionType:
		consumed := c.onKeyboardTypeEvent(element, event)
		if consumed {
			element.Invalidate()
		}
		return consumed

	default:
		return false
	}
}

func (c *editBoxComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	switch event.Action {
	case ui.MouseActionScroll:
		c.offsetX -= event.ScrollX
		c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
		element.Invalidate()
		return true

	case ui.MouseActionDown:
		if event.Button != ui.MouseButtonLeft {
			return false
		}
		c.isDragging = true
		c.cursorColumn = c.findCursorColumn(element, event.X)
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		element.Invalidate()
		return true

	case ui.MouseActionMove: // TODO: Use dragging event
		if c.isDragging {
			c.cursorColumn = c.findCursorColumn(element, event.X)
			element.Invalidate()
		}
		return true

	case ui.MouseActionUp:
		if event.Button != ui.MouseButtonLeft {
			return false
		}
		if c.isDragging {
			c.isDragging = false
			c.cursorColumn = c.findCursorColumn(element, event.X)
			element.Invalidate()
		}
		return true

	default:
		return false
	}
}

func (c *editBoxComponent) OnUndo(element *ui.Element) bool {
	canUndo := c.history.CanUndo()
	if canUndo {
		c.history.Undo()
		c.notifyChanged()
	}
	return canUndo
}

func (c *editBoxComponent) OnRedo(element *ui.Element) bool {
	canRedo := c.history.CanRedo()
	if canRedo {
		c.history.Redo()
		c.notifyChanged()
	}
	return canRedo
}

func (c *editBoxComponent) OnClipboardEvent(element *ui.Element, event ui.ClipboardEvent) bool {
	switch event.Action {
	case ui.ClipboardActionCut:
		if c.isReadOnly {
			return false
		}
		if c.hasSelection() {
			text := string(c.selectedText())
			element.Window().RequestCopy(text)
			c.history.Do(c.changeDeleteSelection())
			c.notifyChanged()
		}
		return true

	case ui.ClipboardActionCopy:
		if c.hasSelection() {
			text := string(c.selectedText())
			element.Window().RequestCopy(text)
		}
		return true

	case ui.ClipboardActionPaste:
		if c.isReadOnly {
			return false
		}
		if c.hasSelection() {
			c.history.Do(c.changeReplaceSelection([]rune(event.Text)))
		} else {
			c.history.Do(c.changeAppendText([]rune(event.Text)))
		}
		c.notifyChanged()
		return true

	default:
		return false
	}
}

func (c *editBoxComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	os := element.Window().Platform().OS()
	if shortcuts.IsCut(os, event) {
		element.Window().Cut()
		return true
	}
	if shortcuts.IsCopy(os, event) {
		element.Window().Copy()
		return true
	}
	if shortcuts.IsPaste(os, event) {
		element.Window().Paste()
		return true
	}
	if shortcuts.IsUndo(os, event) {
		element.Window().Undo()
		return true
	}
	if shortcuts.IsRedo(os, event) {
		element.Window().Redo()
		return true
	}
	if shortcuts.IsSelectAll(os, event) {
		c.selectAll()
		return true
	}

	switch event.Code {

	case ui.KeyCodeEscape:
		c.isDragging = false
		c.resetSelector()
		return true

	case ui.KeyCodeHome:
		if !c.isReadOnly {
			c.moveCursorToStartOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeEnd:
		if !c.isReadOnly {
			c.moveCursorToEndOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeArrowUp:
		if !c.isReadOnly {
			c.moveCursorToStartOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeArrowDown:
		if !c.isReadOnly {
			c.moveCursorToEndOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeArrowLeft:
		if c.isReadOnly {
			c.scrollLeft()
		} else {
			if event.Modifiers.Contains(ui.KeyModifierShift) {
				c.moveCursorLeft()
			} else {
				if c.hasSelection() {
					c.moveCursorToSelectionStart()
				} else {
					c.moveCursorLeft()
				}
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeArrowRight:
		if c.isReadOnly {
			c.scrollRight()
		} else {
			if event.Modifiers.Contains(ui.KeyModifierShift) {
				c.moveCursorRight()
			} else {
				if c.hasSelection() {
					c.moveCursorToSelectionEnd()
				} else {
					c.moveCursorRight()
				}
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeBackspace:
		if c.isReadOnly {
			return false
		}
		if c.hasSelection() {
			c.history.Do(c.changeDeleteSelection())
		} else {
			c.history.Do(c.changeDeleteCharacterLeft())
		}
		c.notifyChanged()
		return true

	case ui.KeyCodeDelete:
		if c.isReadOnly {
			return false
		}
		if c.hasSelection() {
			c.history.Do(c.changeDeleteSelection())
		} else {
			c.history.Do(c.changeDeleteCharacterRight())
		}
		c.notifyChanged()
		return true

	case ui.KeyCodeEnter:
		if c.isReadOnly {
			return false
		}
		c.notifySubmitted()
		return true

	case ui.KeyCodeTab:
		event.Rune = '\t'
		return c.onKeyboardTypeEvent(element, event)

	default:
		return false
	}
}

func (c *editBoxComponent) onKeyboardTypeEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if c.isReadOnly {
		return false
	}
	if c.hasSelection() {
		c.history.Do(c.changeReplaceSelection([]rune{event.Rune}))
	} else {
		c.history.Do(c.changeAppendText([]rune{event.Rune}))
	}
	c.notifyChanged()
	return true
}

func (c *editBoxComponent) hasSelection() bool {
	return c.cursorColumn != c.selectorColumn
}

func (c *editBoxComponent) selectedText() []rune {
	fromColumn := min(c.cursorColumn, c.selectorColumn)
	toColumn := max(c.cursorColumn, c.selectorColumn)
	return slices.Clone(c.line[fromColumn:toColumn])
}

func (c *editBoxComponent) selectAll() {
	c.selectorColumn = 0
	c.cursorColumn = len(c.line)
}

func (c *editBoxComponent) scrollLeft() {
	c.offsetX -= 20
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editBoxComponent) scrollRight() {
	c.offsetX += 20
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editBoxComponent) resetSelector() {
	c.selectorColumn = c.cursorColumn
}

func (c *editBoxComponent) moveCursorLeft() {
	if c.cursorColumn > 0 {
		c.cursorColumn--
	}
}

func (c *editBoxComponent) moveCursorRight() {
	if c.cursorColumn < len(c.line) {
		c.cursorColumn++
	}
}

func (c *editBoxComponent) moveCursorToSelectionStart() {
	c.cursorColumn = min(c.cursorColumn, c.selectorColumn)
}

func (c *editBoxComponent) moveCursorToSelectionEnd() {
	c.cursorColumn = max(c.cursorColumn, c.selectorColumn)
}

func (c *editBoxComponent) moveCursorToStartOfLine() {
	c.cursorColumn = 0
}

func (c *editBoxComponent) moveCursorToEndOfLine() {
	c.cursorColumn = len(c.line)
}

func (c *editBoxComponent) changeAppendText(text []rune) state.Change {
	lng := len(text)
	return &editboxChange{
		when: time.Now(),
		forward: []func(){
			c.actionInsertText(c.cursorColumn, text),
			c.actionRelocateCursor(c.cursorColumn + lng),
			c.actionRelocateSelector(c.cursorColumn + lng),
		},
		reverse: []func(){
			c.actionRelocateSelector(c.selectorColumn),
			c.actionRelocateCursor(c.cursorColumn),
			c.actionDeleteText(c.cursorColumn, c.cursorColumn+lng),
		},
	}
}

func (c *editBoxComponent) changeReplaceSelection(text []rune) state.Change {
	fromColumn := min(c.cursorColumn, c.selectorColumn)
	toColumn := max(c.cursorColumn, c.selectorColumn)
	selectedText := slices.Clone(c.line[fromColumn:toColumn])
	return &editboxChange{
		when: time.Now(),
		forward: []func(){
			c.actionDeleteText(fromColumn, toColumn),
			c.actionInsertText(fromColumn, text),
			c.actionRelocateCursor(fromColumn + len(text)),
			c.actionRelocateSelector(fromColumn + len(text)),
		},
		reverse: []func(){
			c.actionRelocateCursor(c.cursorColumn),
			c.actionRelocateSelector(c.selectorColumn),
			c.actionDeleteText(fromColumn, fromColumn+len(text)),
			c.actionInsertText(fromColumn, selectedText),
		},
	}
}

func (c *editBoxComponent) changeDeleteSelection() state.Change {
	fromColumn := min(c.cursorColumn, c.selectorColumn)
	toColumn := max(c.cursorColumn, c.selectorColumn)
	selectedText := slices.Clone(c.line[fromColumn:toColumn])
	return &editboxChange{
		when: time.Now(),
		forward: []func(){
			c.actionRelocateSelector(fromColumn),
			c.actionRelocateCursor(fromColumn),
			c.actionDeleteText(fromColumn, toColumn),
		},
		reverse: []func(){
			c.actionInsertText(fromColumn, selectedText),
			c.actionRelocateCursor(c.cursorColumn),
			c.actionRelocateSelector(c.selectorColumn),
		},
	}
}

func (c *editBoxComponent) changeDeleteCharacterLeft() state.Change {
	if c.cursorColumn <= 0 {
		return emptyEditBoxChange()
	}
	deletedCharacter := c.line[c.cursorColumn-1]
	return &editboxChange{
		when: time.Now(),
		forward: []func(){
			c.actionRelocateCursor(c.cursorColumn - 1),
			c.actionRelocateSelector(c.cursorColumn - 1),
			c.actionDeleteText(c.cursorColumn-1, c.cursorColumn),
		},
		reverse: []func(){
			c.actionInsertText(c.cursorColumn-1, []rune{deletedCharacter}),
			c.actionRelocateSelector(c.selectorColumn),
			c.actionRelocateCursor(c.cursorColumn),
		},
	}
}

func (c *editBoxComponent) changeDeleteCharacterRight() state.Change {
	if c.cursorColumn >= len(c.line) {
		return emptyEditBoxChange()
	}
	deletedCharacter := c.line[c.cursorColumn]
	return &editboxChange{
		when: time.Now(),
		forward: []func(){
			c.actionRelocateCursor(c.cursorColumn),
			c.actionRelocateSelector(c.cursorColumn),
			c.actionDeleteText(c.cursorColumn, c.cursorColumn+1),
		},
		reverse: []func(){
			c.actionInsertText(c.cursorColumn, []rune{deletedCharacter}),
			c.actionRelocateSelector(c.selectorColumn),
			c.actionRelocateCursor(c.cursorColumn),
		},
	}
}

func (c *editBoxComponent) actionInsertText(position int, text []rune) func() {
	return func() {
		preText := c.line[:position]
		postText := c.line[position:]
		c.line = gog.Concat(
			preText,
			text,
			postText,
		)
	}
}

func (c *editBoxComponent) actionDeleteText(fromPosition, toPosition int) func() {
	return func() {
		c.line = slices.Delete(c.line, fromPosition, toPosition)
	}
}

func (c *editBoxComponent) actionRelocateCursor(position int) func() {
	return func() {
		c.cursorColumn = position
	}
}

func (c *editBoxComponent) actionRelocateSelector(position int) func() {
	return func() {
		c.selectorColumn = position
	}
}

func (c *editBoxComponent) findCursorColumn(element *ui.Element, x int) int {
	x -= element.Padding().Left - int(c.offsetX)

	bestColumn := 0
	bestDistance := abs(x)

	column := 1
	offset := float32(0.0)
	iterator := c.font.TextIterator(string(c.line), c.fontSize)
	for iterator.Next() {
		character := iterator.Character()
		offset += character.Kern + character.Width
		if distance := abs(x - int(offset)); distance < bestDistance {
			bestColumn = column
			bestDistance = distance
		}
		column++
	}
	return bestColumn
}

func (c *editBoxComponent) notifyChanged() {
	if c.onChange != nil {
		c.onChange(string(c.line))
	}
}

func (c *editBoxComponent) notifySubmitted() {
	if c.onSubmit != nil {
		c.onSubmit(string(c.line))
	}
}

func emptyEditBoxChange() *editboxChange {
	return &editboxChange{
		when: time.Now(),
	}
}

type editboxChange struct {
	when    time.Time
	forward []func()
	reverse []func()
}

func (c *editboxChange) Apply() {
	for _, action := range c.forward {
		action()
	}
}

func (c *editboxChange) Revert() {
	for _, action := range c.reverse {
		action()
	}
}

func (c *editboxChange) Extend(other state.Change) bool {
	otherChange, ok := other.(*editboxChange)
	if !ok {
		return false
	}
	if otherChange.when.Sub(c.when) > 500*time.Millisecond {
		return false
	}
	c.forward = append(c.forward, otherChange.forward...)
	c.reverse = append(otherChange.reverse, c.reverse...)
	c.when = otherChange.when
	return true
}

func abs(a int) int {
	return max(a, -a)
}
