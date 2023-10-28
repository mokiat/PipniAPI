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
	editboxHistoryCapacity  = 100
	editboxPaddingLeft      = 10
	editboxPaddingRight     = 10
	editboxPaddingTop       = 5
	editboxPaddingBottom    = 5
	editboxTextPaddingLeft  = 2
	editboxTextPaddingRight = 2
	editboxCursorWidth      = float32(1.0)
	editboxBorderSize       = float32(2.0)
	editboxBorderRadius     = float32(8.0)
	editboxKeyScrollSpeed   = 20
	editboxFontSize         = float32(18.0)
)

var EditBox = co.Define(&editboxComponent{})

type EditBoxData struct {
	ReadOnly bool
	Text     string
}

type EditBoxCallbackData struct {
	OnChange func(string)
	OnSubmit func(string)
}

var _ ui.ElementHistoryHandler = (*editboxComponent)(nil)
var _ ui.ElementClipboardHandler = (*editboxComponent)(nil)
var _ ui.ElementResizeHandler = (*editboxComponent)(nil)
var _ ui.ElementRenderHandler = (*editboxComponent)(nil)
var _ ui.ElementKeyboardHandler = (*editboxComponent)(nil)
var _ ui.ElementMouseHandler = (*editboxComponent)(nil)

type editboxComponent struct {
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

	textWidth  int
	textHeight int

	offsetX    float32
	maxOffsetX float32

	isDragging bool
}

func (c *editboxComponent) OnCreate() {
	c.history = state.NewHistory(editboxHistoryCapacity)

	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = editboxFontSize

	c.cursorColumn = 0
	c.selectorColumn = 0

	data := co.GetData[EditBoxData](c.Properties())
	c.isReadOnly = data.ReadOnly
	c.line = []rune(data.Text)
	c.refreshTextSize()
}

func (c *editboxComponent) OnUpsert() {
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

func (c *editboxComponent) Render() co.Instance {
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
				Width:  c.textWidth + textPadding,
				Height: c.textHeight,
			}.Grow(padding.Size())),
			Padding: padding,
		})
	})
}

func (c *editboxComponent) OnUndo(element *ui.Element) bool {
	canUndo := c.history.CanUndo()
	if canUndo {
		c.history.Undo()
		c.notifyChanged()
	}
	return canUndo
}

func (c *editboxComponent) OnRedo(element *ui.Element) bool {
	canRedo := c.history.CanRedo()
	if canRedo {
		c.history.Redo()
		c.notifyChanged()
	}
	return canRedo
}

func (c *editboxComponent) OnClipboardEvent(element *ui.Element, event ui.ClipboardEvent) bool {
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

func (c *editboxComponent) OnResize(element *ui.Element, bounds ui.Bounds) {
	c.refreshScrollBounds(element)
}

func (c *editboxComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	c.refreshScrollBounds(element)

	bounds := canvas.DrawBounds(element, false)
	contentBounds := canvas.DrawBounds(element, true)

	canvas.Push()
	canvas.ClipRect(bounds.Position, bounds.Size)
	canvas.Translate(bounds.Position)
	c.drawFrame(element, canvas, bounds.Size)
	canvas.Pop()

	canvas.Push()
	canvas.ClipRect(contentBounds.Position, contentBounds.Size)
	canvas.Translate(contentBounds.Position)
	c.drawContent(element, canvas, contentBounds.Size)
	canvas.Pop()

	canvas.Push()
	canvas.ClipRect(bounds.Position, bounds.Size)
	canvas.Translate(bounds.Position)
	c.drawFrameBorder(element, canvas, bounds.Size)
	canvas.Pop()
}

func (c *editboxComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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

func (c *editboxComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
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

	case ui.MouseActionMove:
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

func (c *editboxComponent) drawFrame(element *ui.Element, canvas *ui.Canvas, bounds sprec.Vec2) {
	canvas.Reset()
	canvas.RoundRectangle(
		sprec.ZeroVec2(),
		bounds,
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

func (c *editboxComponent) drawFrameBorder(element *ui.Element, canvas *ui.Canvas, bounds sprec.Vec2) {
	canvas.Reset()
	if element.IsFocused() {
		canvas.SetStrokeColor(std.SecondaryLightColor)
	} else {
		canvas.SetStrokeColor(std.PrimaryLightColor)
	}
	canvas.SetStrokeSize(editboxBorderSize)
	canvas.RoundRectangle(
		sprec.ZeroVec2(),
		bounds,
		sprec.NewVec4(
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
			editboxBorderRadius,
		),
	)
	canvas.Stroke()
}

func (c *editboxComponent) drawContent(element *ui.Element, canvas *ui.Canvas, bounds sprec.Vec2) {
	c.drawSelection(element, canvas)
	c.drawText(element, canvas, bounds)
	c.drawCursor(element, canvas)
}

func (c *editboxComponent) drawSelection(element *ui.Element, canvas *ui.Canvas) {
	if !c.hasSelection() || !element.IsFocused() {
		return
	}

	fromColumn, toColumn := c.selectionRange()
	preSelectionWidth := c.font.LineWidth(c.line[:fromColumn], c.fontSize)
	selectionWidth := c.font.LineWidth(c.line[fromColumn:toColumn], c.fontSize)
	selectionHeight := c.font.LineHeight(c.fontSize)

	selectionPosition := sprec.Vec2{
		X: editboxTextPaddingLeft + preSelectionWidth - c.offsetX,
	}
	selectionSize := sprec.Vec2{
		X: selectionWidth,
		Y: selectionHeight,
	}
	canvas.Reset()
	canvas.Rectangle(selectionPosition, selectionSize)
	canvas.Fill(ui.Fill{
		Color: std.SecondaryLightColor,
	})
}

func (c *editboxComponent) drawText(element *ui.Element, canvas *ui.Canvas, bounds sprec.Vec2) {
	if len(c.line) == 0 {
		return
	}
	fromColumn, toColumn := c.visibleColumns(bounds)
	if fromColumn > toColumn {
		return
	}

	preVisibleText := c.line[:fromColumn]
	preVisibleTextWidth := c.font.LineWidth(preVisibleText, c.fontSize)
	visibleText := c.line[fromColumn : toColumn+1]
	visibleTextPosition := sprec.Vec2{
		X: editboxTextPaddingLeft + preVisibleTextWidth - c.offsetX,
	}

	canvas.Reset()
	canvas.FillText(string(visibleText), visibleTextPosition, ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: std.OnSurfaceColor,
	})
}

func (c *editboxComponent) drawCursor(element *ui.Element, canvas *ui.Canvas) {
	if c.isReadOnly || !element.IsFocused() {
		return
	}

	preCursorText := c.line[:c.cursorColumn]
	preCursorTextWidth := c.font.LineWidth(preCursorText, c.fontSize)
	preCursorTextHeight := c.font.LineHeight(c.fontSize)

	cursorPosition := sprec.Vec2{
		X: editboxTextPaddingLeft + preCursorTextWidth - c.offsetX,
	}
	cursorSize := sprec.Vec2{
		X: editboxCursorWidth,
		Y: preCursorTextHeight,
	}
	canvas.Reset()
	canvas.Rectangle(cursorPosition, cursorSize)
	canvas.Fill(ui.Fill{
		Color: std.PrimaryColor,
	})
}

func (c *editboxComponent) visibleColumns(bounds sprec.Vec2) (int, int) {
	minVisible := len(c.line)
	maxVisible := -1
	offset := float32(editboxTextPaddingLeft) - float32(c.offsetX)
	iterator := c.font.LineIterator(c.line, c.fontSize)
	column := 0
	for iterator.Next() {
		character := iterator.Character()
		characterWidth := character.Kern + character.Width
		if offset+characterWidth > 0.0 && offset < bounds.X {
			minVisible = min(minVisible, column)
			maxVisible = max(maxVisible, column)
		}
		offset += characterWidth
		column++
	}
	return minVisible, maxVisible
}

func (c *editboxComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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

func (c *editboxComponent) onKeyboardTypeEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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

func (c *editboxComponent) hasSelection() bool {
	return c.cursorColumn != c.selectorColumn
}

func (c *editboxComponent) selectionRange() (int, int) {
	fromColumn := min(c.cursorColumn, c.selectorColumn)
	toColumn := max(c.cursorColumn, c.selectorColumn)
	return fromColumn, toColumn
}

func (c *editboxComponent) selectedText() []rune {
	fromColumn, toColumn := c.selectionRange()
	return slices.Clone(c.line[fromColumn:toColumn])
}

func (c *editboxComponent) selectAll() {
	c.selectorColumn = 0
	c.cursorColumn = len(c.line)
}

func (c *editboxComponent) scrollLeft() {
	c.offsetX -= editboxKeyScrollSpeed
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editboxComponent) scrollRight() {
	c.offsetX += editboxKeyScrollSpeed
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editboxComponent) resetSelector() {
	c.selectorColumn = c.cursorColumn
}

func (c *editboxComponent) moveCursorLeft() {
	if c.cursorColumn > 0 {
		c.cursorColumn--
	}
}

func (c *editboxComponent) moveCursorRight() {
	if c.cursorColumn < len(c.line) {
		c.cursorColumn++
	}
}

func (c *editboxComponent) moveCursorToSelectionStart() {
	c.cursorColumn = min(c.cursorColumn, c.selectorColumn)
}

func (c *editboxComponent) moveCursorToSelectionEnd() {
	c.cursorColumn = max(c.cursorColumn, c.selectorColumn)
}

func (c *editboxComponent) moveCursorToStartOfLine() {
	c.cursorColumn = 0
}

func (c *editboxComponent) moveCursorToEndOfLine() {
	c.cursorColumn = len(c.line)
}

func (c *editboxComponent) changeAppendText(text []rune) state.Change {
	lng := len(text)
	return &textTypeChange{
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

func (c *editboxComponent) changeReplaceSelection(text []rune) state.Change {
	fromColumn, toColumn := c.selectionRange()
	selectedText := slices.Clone(c.line[fromColumn:toColumn])
	return &textTypeChange{
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

func (c *editboxComponent) changeDeleteSelection() state.Change {
	fromColumn, toColumn := c.selectionRange()
	selectedText := slices.Clone(c.line[fromColumn:toColumn])
	return &textTypeChange{
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

func (c *editboxComponent) changeDeleteCharacterLeft() state.Change {
	if c.cursorColumn == 0 {
		return emptyTextTypeChange()
	}
	deletedCharacter := c.line[c.cursorColumn-1]
	return &textTypeChange{
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

func (c *editboxComponent) changeDeleteCharacterRight() state.Change {
	if c.cursorColumn >= len(c.line) {
		return emptyTextTypeChange()
	}
	deletedCharacter := c.line[c.cursorColumn]
	return &textTypeChange{
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

func (c *editboxComponent) actionInsertText(position int, text []rune) func() {
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

func (c *editboxComponent) actionDeleteText(fromPosition, toPosition int) func() {
	return func() {
		c.line = slices.Delete(c.line, fromPosition, toPosition)
	}
}

func (c *editboxComponent) actionRelocateCursor(position int) func() {
	return func() {
		c.cursorColumn = position
	}
}

func (c *editboxComponent) actionRelocateSelector(position int) func() {
	return func() {
		c.selectorColumn = position
	}
}

func (c *editboxComponent) findCursorColumn(element *ui.Element, x int) int {
	x -= element.Padding().Left - int(c.offsetX) + editboxTextPaddingLeft

	bestColumn := 0
	bestDistance := abs(x)

	column := 1
	offset := float32(0.0)
	iterator := c.font.LineIterator(c.line, c.fontSize)
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

func (c *editboxComponent) refreshTextSize() {
	c.textWidth = editboxTextPaddingLeft + int(c.font.LineWidth(c.line, c.fontSize)) + editboxTextPaddingRight
	c.textHeight = int(c.font.LineHeight(c.fontSize))
}

func (c *editboxComponent) refreshScrollBounds(element *ui.Element) {
	bounds := element.ContentBounds()
	availableTextWidth := bounds.Width - editboxTextPaddingLeft - editboxTextPaddingRight
	c.maxOffsetX = float32(max(c.textWidth-availableTextWidth, 0))
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *editboxComponent) notifyChanged() {
	c.refreshTextSize()
	if c.onChange != nil {
		c.onChange(string(c.line))
	}
}

func (c *editboxComponent) notifySubmitted() {
	if c.onSubmit != nil {
		c.onSubmit(string(c.line))
	}
}

func abs(a int) int {
	return max(a, -a)
}
