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
	editBoxHistoryCapacity  = 100
	editboxPaddingLeft      = 10
	editboxPaddingRight     = 10
	editboxPaddingTop       = 5
	editboxPaddingBottom    = 5
	editboxTextPaddingLeft  = 2
	editboxTextPaddingRight = 2
	editboxCursorWidth      = float32(1.0)
	editboxBorderSize       = float32(2.0)
	editboxBorderRadius     = float32(8.0)
	editboxFontSize         = float32(18.0)
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

var _ ui.ElementHistoryHandler = (*editBoxComponent)(nil)
var _ ui.ElementClipboardHandler = (*editBoxComponent)(nil)
var _ ui.ElementResizeHandler = (*editBoxComponent)(nil)
var _ ui.ElementRenderHandler = (*editBoxComponent)(nil)
var _ ui.ElementKeyboardHandler = (*editBoxComponent)(nil)
var _ ui.ElementMouseHandler = (*editBoxComponent)(nil)

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

	textHeight int
	textWidth  int

	offsetX    int
	maxOffsetX int

	isDragging bool
}

func (c *editBoxComponent) OnCreate() {
	c.history = state.NewHistory(editBoxHistoryCapacity)

	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = editboxFontSize

	c.cursorColumn = 0
	c.selectorColumn = 0

	data := co.GetData[EditBoxData](c.Properties())
	c.isReadOnly = data.ReadOnly
	c.line = []rune(data.Text)
	c.refreshTextSize()
}

func (c *editBoxComponent) OnUpsert() {
	data := co.GetData[EditBoxData](c.Properties())
	if data.ReadOnly != c.isReadOnly {
		c.isReadOnly = data.ReadOnly
		c.history.Clear()
	}
	if data.Text != string(c.line) {
		c.history.Clear()
		c.line = []rune(data.Text)
		c.refreshTextSize()
	}

	callbackData := co.GetOptionalCallbackData(c.Properties(), EditBoxCallbackData{})
	c.onChange = callbackData.OnChange
	c.onSubmit = callbackData.OnSubmit

	c.cursorColumn = min(c.cursorColumn, len(c.line))
	c.selectorColumn = min(c.selectorColumn, len(c.line))
}

func (c *editBoxComponent) Render() co.Instance {
	padding := ui.Spacing{
		Left:   editboxPaddingLeft,
		Right:  editboxPaddingRight,
		Top:    editboxPaddingTop,
		Bottom: editboxPaddingBottom,
	}
	textPadding := editboxTextPaddingLeft + editboxTextPaddingRight

	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			IdealSize: opt.V(ui.Size{
				Width:  c.textWidth + padding.Horizontal() + textPadding,
				Height: c.textHeight + padding.Vertical(),
			}),
			Padding: padding,
		})
	})
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

func (c *editBoxComponent) OnResize(element *ui.Element, bounds ui.Bounds) {
	c.refreshScrollBounds(element)
}

func (c *editBoxComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	c.refreshScrollBounds(element)
	c.drawFrame(element, canvas)
	c.drawContent(element, canvas)
	c.drawFrameBorder(element, canvas)
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
		if event.Modifiers.Contains(ui.KeyModifierShift) && (event.ScrollX == 0) {
			c.offsetX -= event.ScrollY
		} else {
			c.offsetX -= event.ScrollX
		}
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

func (c *editBoxComponent) drawFrame(element *ui.Element, canvas *ui.Canvas) {
	bounds := canvas.DrawBounds(element, false)

	canvas.Reset()
	canvas.RoundRectangle(
		bounds.Position,
		bounds.Size,
		sprec.NewVec4(
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
		),
	)
	canvas.Fill(ui.Fill{
		Color: std.SurfaceColor,
	})
}

func (c *editBoxComponent) drawFrameBorder(element *ui.Element, canvas *ui.Canvas) {
	bounds := canvas.DrawBounds(element, false)

	canvas.Reset()
	if element.IsFocused() {
		canvas.SetStrokeColor(std.SecondaryLightColor)
	} else {
		canvas.SetStrokeColor(std.PrimaryLightColor)
	}
	canvas.SetStrokeSize(editboxBorderSize)
	canvas.RoundRectangle(
		bounds.Position,
		bounds.Size,
		sprec.NewVec4(
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
		),
	)
	canvas.Stroke()
}

func (c *editBoxComponent) drawContent(element *ui.Element, canvas *ui.Canvas) {
	contentBounds := canvas.DrawBounds(element, true)

	canvas.Push()
	canvas.SetClipRect(
		contentBounds.X(),
		contentBounds.X()+contentBounds.Width(),
		contentBounds.Y(),
		contentBounds.Y()+contentBounds.Height(),
	)
	canvas.Translate(sprec.Vec2{
		X: contentBounds.X() + float32(editboxTextPaddingLeft) - float32(c.offsetX),
		Y: contentBounds.Y(),
	})
	c.drawSelection(element, canvas)
	c.drawText(element, canvas)
	c.drawCursor(element, canvas)
	canvas.Pop()
}

func (c *editBoxComponent) drawSelection(element *ui.Element, canvas *ui.Canvas) {
	if !c.hasSelection() || !element.IsFocused() {
		return
	}

	fromColumn, toColumn := c.selectionRange()
	preSelectionWidth := c.font.LineWidth(c.line[:fromColumn], c.fontSize)
	selectionWidth := c.font.LineWidth(c.line[fromColumn:toColumn], c.fontSize)
	selectionHeight := c.font.LineHeight(c.fontSize)

	selectionPosition := sprec.NewVec2(preSelectionWidth, 0.0)
	selectionSize := sprec.NewVec2(selectionWidth, selectionHeight)

	canvas.Reset()
	canvas.Rectangle(selectionPosition, selectionSize)
	canvas.Fill(ui.Fill{
		Color: std.SecondaryLightColor,
	})
}

func (c *editBoxComponent) drawText(element *ui.Element, canvas *ui.Canvas) {
	if len(c.line) == 0 {
		return
	}

	canvas.Reset()
	canvas.FillText(string(c.line), sprec.ZeroVec2(), ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: std.OnSurfaceColor,
	})
}

func (c *editBoxComponent) drawCursor(element *ui.Element, canvas *ui.Canvas) {
	if c.isReadOnly || !element.IsFocused() {
		return
	}

	preCursorText := c.line[:c.cursorColumn]
	preCursorTextWidth := c.font.LineWidth(preCursorText, c.fontSize)
	preCursorTextHeight := c.font.LineHeight(c.fontSize)

	cursorPosition := sprec.NewVec2(preCursorTextWidth, 0.0)
	cursorSize := sprec.NewVec2(editboxCursorWidth, preCursorTextHeight)

	canvas.Reset()
	canvas.Rectangle(cursorPosition, cursorSize)
	canvas.Fill(ui.Fill{
		Color: std.PrimaryColor,
	})
}

func (c *editBoxComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	os := element.Window().Platform().OS()
	if shortcuts.IsClose(os, event) {
		return false // propagate up
	}
	if shortcuts.IsSave(os, event) {
		return false // propagate up
	}
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
	if shortcuts.IsJumpToLineStart(os, event) {
		if !c.isReadOnly {
			c.moveCursorToStartOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true
	}
	if shortcuts.IsJumpToLineEnd(os, event) {
		if !c.isReadOnly {
			c.moveCursorToEndOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true
	}

	switch event.Code {

	case ui.KeyCodeEscape:
		c.isDragging = false
		c.resetSelector()
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
		return false

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

func (c *editBoxComponent) selectionRange() (int, int) {
	fromColumn := min(c.cursorColumn, c.selectorColumn)
	toColumn := max(c.cursorColumn, c.selectorColumn)
	return fromColumn, toColumn
}

func (c *editBoxComponent) selectedText() []rune {
	fromColumn, toColumn := c.selectionRange()
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

func (c *editBoxComponent) refreshTextSize() {
	c.textWidth = int(c.font.LineWidth(c.line, c.fontSize))
	c.textHeight = int(c.font.LineHeight(c.fontSize))
}

func (c *editBoxComponent) refreshScrollBounds(element *ui.Element) {
	bounds := element.ContentBounds()
	availableTextWidth := bounds.Width - editboxTextPaddingLeft - editboxTextPaddingRight
	c.maxOffsetX = max(c.textWidth-availableTextWidth, 0)
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editBoxComponent) notifyChanged() {
	c.refreshTextSize()
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
