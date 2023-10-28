package widget

import (
	"math"
	"strings"

	"github.com/mokiat/PipniAPI/internal/shortcuts"
	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/state"
	"github.com/mokiat/lacking/ui/std"
	"golang.org/x/exp/slices"
)

const (
	codeAreaHistoryCapacity   = 100
	codeAreaPaddingLeft       = 2
	codeAreaPaddingRight      = 2
	codeAreaPaddingTop        = 2
	codeAreaPaddingBottom     = 2
	codeAreaTextPaddingLeft   = 5
	codeAreaTextPaddingRight  = 5
	codeAreaRulerPaddingLeft  = 5
	codeAreaRulerPaddingRight = 5
	codeAreaCursorWidth       = float32(1.0)
	codeAreaBorderSize        = float32(2.0)
	codeAreaKeyScrollSpeed    = 20
	codeAreaFontSize          = float32(18.0)
)

var CodeArea = co.Define(&codeAreaComponent{})

type CodeAreaData struct {
	ReadOnly bool
	Code     string
}

type CodeAreaCallbackData struct {
	OnChange func(string)
}

var _ ui.ElementHistoryHandler = (*codeAreaComponent)(nil)
var _ ui.ElementClipboardHandler = (*codeAreaComponent)(nil)
var _ ui.ElementResizeHandler = (*codeAreaComponent)(nil)
var _ ui.ElementRenderHandler = (*codeAreaComponent)(nil)
var _ ui.ElementKeyboardHandler = (*codeAreaComponent)(nil)
var _ ui.ElementMouseHandler = (*codeAreaComponent)(nil)

type codeAreaComponent struct {
	co.BaseComponent

	history *state.History

	font     *ui.Font
	fontSize float32

	cursorRow      int
	cursorColumn   int
	selectorRow    int
	selectorColumn int

	isReadOnly bool
	lines      [][]rune
	onChange   func(string)

	textWidth  int
	textHeight int
	rulerWidth int

	offsetX    float32
	offsetY    float32
	maxOffsetX float32
	maxOffsetY float32

	isDragging bool
}

func (c *codeAreaComponent) OnCreate() {
	c.history = state.NewHistory(editboxHistoryCapacity)

	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = codeAreaFontSize

	c.cursorRow = 0
	c.cursorColumn = 0

	data := co.GetData[CodeAreaData](c.Properties())
	c.isReadOnly = data.ReadOnly
	c.lines = splitLines(data.Code)
	c.refreshTextSize()
}

func (c *codeAreaComponent) OnUpsert() {
	data := co.GetData[CodeAreaData](c.Properties())
	if data.ReadOnly != c.isReadOnly {
		c.isReadOnly = data.ReadOnly
		c.history.Clear()
	}
	if data.Code != c.constructText() {
		c.history.Clear()
		c.lines = splitLines(data.Code)
		c.refreshTextSize()
	}

	callbackData := co.GetOptionalCallbackData[CodeAreaCallbackData](c.Properties(), CodeAreaCallbackData{})
	c.onChange = callbackData.OnChange

	c.cursorRow = min(c.cursorRow, len(c.lines)-1)
	c.cursorColumn = min(c.cursorColumn, len(c.lines[c.cursorRow]))
	c.selectorRow = min(c.selectorRow, len(c.lines)-1)
	c.selectorColumn = min(c.selectorColumn, len(c.lines[c.selectorRow]))
}

func (c *codeAreaComponent) Render() co.Instance {
	padding := ui.Spacing{
		Left:   codeAreaPaddingLeft,
		Right:  codeAreaPaddingRight,
		Top:    codeAreaPaddingTop,
		Bottom: codeAreaPaddingBottom,
	}
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			IdealSize: opt.V(ui.Size{
				Width:  c.textWidth + c.rulerWidth,
				Height: c.textHeight,
			}.Grow(padding.Size())),
		})
	})
}

func (c *codeAreaComponent) OnUndo(element *ui.Element) bool {
	canUndo := c.history.CanUndo()
	if canUndo {
		c.history.Undo()
		c.notifyChanged()
	}
	return canUndo
}

func (c *codeAreaComponent) OnRedo(element *ui.Element) bool {
	canRedo := c.history.CanRedo()
	if canRedo {
		c.history.Redo()
		c.notifyChanged()
	}
	return canRedo
}

func (c *codeAreaComponent) OnClipboardEvent(element *ui.Element, event ui.ClipboardEvent) bool {
	switch event.Action {
	case ui.ClipboardActionCut:
		if c.isReadOnly {
			return false
		}
		if c.hasSelection() {
			// TODO
			// 		text := string(c.selectedText())
			// 		element.Window().RequestCopy(text)
			// 		c.history.Do(c.changeDeleteSelection())
			c.notifyChanged()
		}
		return true

	case ui.ClipboardActionCopy:
		if c.hasSelection() {
			text := strings.Join(gog.Map(c.selectedLines(), lineToText), "\n")
			element.Window().RequestCopy(text)
		}
		return true

	case ui.ClipboardActionPaste:
		if c.isReadOnly {
			return false
		}

		lines := splitLines(event.Text)
		if c.hasSelection() {
			// TODO
			// 		c.history.Do(c.changeReplaceSelection([]rune(event.Text)))
		} else {
			c.history.Do(c.changeAppendText(lines))
		}
		c.notifyChanged()
		return true

	default:
		return false
	}
}

func (c *codeAreaComponent) OnResize(element *ui.Element, bounds ui.Bounds) {
	c.refreshScrollBounds(element)
}

func (c *codeAreaComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
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

func (c *codeAreaComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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

func (c *codeAreaComponent) OnMouseEvent(element *ui.Element, event ui.MouseEvent) bool {
	switch event.Action {
	case ui.MouseActionScroll:
		if event.Modifiers.Contains(ui.KeyModifierShift) && (event.ScrollX == 0) {
			c.offsetX -= event.ScrollY
		} else {
			c.offsetX -= event.ScrollX
			c.offsetY -= event.ScrollY
		}
		c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
		c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
		element.Invalidate()
		return true

	case ui.MouseActionDown:
		if event.Button != ui.MouseButtonLeft {
			return false
		}
		c.isDragging = true
		c.cursorRow = c.findCursorRow(element, event.Y)
		c.cursorColumn = c.findCursorColumn(element, event.X)
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.clearSelection()
		}
		element.Invalidate()
		return true

	case ui.MouseActionMove:
		if c.isDragging {
			c.cursorRow = c.findCursorRow(element, event.Y)
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
			c.cursorRow = c.findCursorRow(element, event.Y)
			c.cursorColumn = c.findCursorColumn(element, event.X)
			element.Invalidate()
		}
		return true

	default:
		return false
	}
}

func (c *codeAreaComponent) refreshTextSize() {
	txtWidth := float32(0.0)
	for _, line := range c.lines {
		lineWidth := c.font.LineWidth(line, c.fontSize)
		txtWidth = max(txtWidth, lineWidth)
	}
	txtHeight := c.font.LineHeight(c.fontSize) * float32(len(c.lines))

	c.textWidth = codeAreaTextPaddingLeft + int(math.Ceil(float64(txtWidth))) + codeAreaTextPaddingRight
	c.textHeight = int(math.Ceil(float64(txtHeight)))

	digitSize := c.font.LineWidth([]rune{'0'}, c.fontSize)
	digitCount := countDigits(len(c.lines))
	rulerTextWidth := int(math.Ceil(float64(digitSize)) * float64(digitCount))
	c.rulerWidth = codeAreaRulerPaddingLeft + rulerTextWidth + codeAreaRulerPaddingRight
}

func (c *codeAreaComponent) refreshScrollBounds(element *ui.Element) {
	bounds := element.ContentBounds()

	textPadding := codeAreaTextPaddingLeft + codeAreaTextPaddingRight
	availableTextWidth := bounds.Width - c.rulerWidth - textPadding
	availableTextHeight := bounds.Height
	c.maxOffsetX = float32(max(c.textWidth-availableTextWidth, 0))
	c.maxOffsetY = float32(max(c.textHeight-availableTextHeight, 0))
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
	c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
}

func (c *codeAreaComponent) findCursorRow(element *ui.Element, y int) int {
	y += int(c.offsetY)
	y -= element.Padding().Top

	lineHeight := c.font.LineHeight(c.fontSize)
	row := y / int(lineHeight)
	return min(max(0, row), len(c.lines)-1)
}

func (c *codeAreaComponent) findCursorColumn(element *ui.Element, x int) int {
	x += int(c.offsetX)
	x -= element.Padding().Left
	x -= c.rulerWidth
	x -= codeAreaTextPaddingLeft

	bestColumn := 0
	bestDistance := abs(x)

	column := 1
	offset := float32(0.0)
	iterator := c.font.LineIterator(c.lines[c.cursorRow], c.fontSize)
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

func (c *codeAreaComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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

	switch event.Code {

	case ui.KeyCodeEscape:
		element.Window().DiscardFocus()
		return true

	case ui.KeyCodeArrowUp:
		if c.isReadOnly {
			c.scrollUp()
		} else {
			c.moveCursorUp()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.clearSelection()
			}
		}
		return true

	case ui.KeyCodeArrowDown:
		if c.isReadOnly {
			c.scrollDown()
		} else {
			c.moveCursorDown()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.clearSelection()
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
				if !c.hasSelection() {
					c.moveCursorLeft()
				}
				c.clearSelection()
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
				if !c.hasSelection() {
					c.moveCursorRight()
				}
				c.clearSelection()
			}
		}
		return true

	case ui.KeyCodeBackspace:
		if c.isReadOnly {
			return false
		}
		// TODO: Check if selection
		c.eraseLeft()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.clearSelection()
		}
		c.notifyChanged()
		return true

	case ui.KeyCodeDelete:
		if c.isReadOnly {
			return false
		}
		// TODO: Check if selection
		c.eraseRight()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.clearSelection()
		}
		c.notifyChanged()
		return true

	case ui.KeyCodeEnter:
		if c.isReadOnly {
			return false
		}
		c.breakLine()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.clearSelection()
		}
		c.notifyChanged()
		return true

	case ui.KeyCodeTab:
		if c.isReadOnly {
			return false
		}
		c.appendCharacter('\t')
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.clearSelection()
		}
		c.notifyChanged()
		return true

	default:
		return false
	}
}

func (c *codeAreaComponent) onKeyboardTypeEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if c.isReadOnly {
		return false
	}
	c.appendCharacter(event.Rune)
	c.clearSelection()
	c.notifyChanged()
	return true
}

func (c *codeAreaComponent) scrollUp() {
	c.offsetY -= 20
	c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
}

func (c *codeAreaComponent) scrollDown() {
	c.offsetY += 20
	c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
}

func (c *codeAreaComponent) scrollLeft() {
	c.offsetX -= 20
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *codeAreaComponent) scrollRight() {
	c.offsetX += 20
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
}

func (c *codeAreaComponent) moveCursorUp() {
	if c.cursorRow > 0 {
		c.cursorRow--
		if c.cursorColumn > len(c.lines[c.cursorRow]) {
			c.cursorColumn = len(c.lines[c.cursorRow])
		}
	} else {
		c.moveCursorToStartOfLine()
	}
}

func (c *codeAreaComponent) moveCursorDown() {
	if c.cursorRow < len(c.lines)-1 {
		c.cursorRow++
		if c.cursorColumn > len(c.lines[c.cursorRow]) {
			c.cursorColumn = len(c.lines[c.cursorRow])
		}
	} else {
		c.moveCursorToEndOfLine()
	}
}

func (c *codeAreaComponent) moveCursorLeft() {
	if c.cursorColumn > 0 {
		c.cursorColumn--
	} else {
		if c.cursorRow > 0 {
			c.moveCursorUp()
			c.moveCursorToEndOfLine()
		}
	}
}

func (c *codeAreaComponent) moveCursorRight() {
	if c.cursorColumn < len(c.lines[c.cursorRow]) {
		c.cursorColumn++
	} else {
		if c.cursorRow < len(c.lines)-1 {
			c.moveCursorDown()
			c.moveCursorToStartOfLine()
		}
	}
}

func (c *codeAreaComponent) moveCursorToStartOfLine() {
	c.cursorColumn = 0
}

func (c *codeAreaComponent) moveCursorToEndOfLine() {
	c.cursorColumn = len(c.lines[c.cursorRow])
}

func (c *codeAreaComponent) changeAppendText(lines [][]rune) state.Change {
	if len(lines) == 0 {
		return emptyTextTypeChange() // TODO: Return nil
	}
	newCursorRow := c.cursorRow + len(lines) - 1
	newCursorColumn := c.cursorColumn + len(lines[0])
	if lng := len(lines); lng > 1 {
		newCursorColumn = len(lines[lng-1])
	}
	forward := []state.Action{
		c.actionInsertText(c.cursorRow, c.cursorColumn, lines[0]),
		c.actionInsertLines(c.cursorRow+1, lines[1:]),
		c.actionRelocateCursor(newCursorRow, newCursorColumn),
		c.actionRelocateSelector(newCursorRow, newCursorColumn),
	}
	reverse := []state.Action{
		c.actionRelocateSelector(c.selectorRow, c.selectorColumn),
		c.actionRelocateCursor(c.cursorRow, c.cursorColumn),
		c.actionDeleteLines(c.cursorRow+1, c.cursorRow+len(lines)),
		c.actionDeleteText(c.cursorRow, c.cursorColumn, c.cursorColumn+len(lines[0])),
	}
	return c.createChange(forward, reverse)
}

func (c *codeAreaComponent) createChange(forward, reverse []state.Action) state.Change {
	return state.AccumActionChange(forward, reverse, textTypeAccumulationDuration)
}

func (c *codeAreaComponent) actionInsertText(row, column int, text []rune) func() {
	return func() {
		line := slices.Clone(c.lines[row])
		preText := line[:column]
		postText := line[column:]
		c.lines[row] = gog.Concat(
			preText,
			text,
			postText,
		)
	}
}

func (c *codeAreaComponent) actionDeleteText(row, fromColumn, toColumn int) func() {
	return func() {
		c.lines[row] = slices.Delete(c.lines[row], fromColumn, toColumn)
	}
}

func (c *codeAreaComponent) actionInsertLines(row int, lines [][]rune) func() {
	return func() {
		c.lines = slices.Insert(c.lines, row, slices.Clone(lines)...)
	}
}

func (c *codeAreaComponent) actionDeleteLines(fromRow, toRow int) func() {
	return func() {
		if fromRow <= toRow {
			c.lines = slices.Delete(c.lines, fromRow, toRow)
		}
	}
}

func (c *codeAreaComponent) actionRelocateCursor(row, column int) func() {
	return func() {
		c.cursorRow = row
		c.cursorColumn = column
	}
}

func (c *codeAreaComponent) actionRelocateSelector(row, column int) func() {
	return func() {
		c.selectorRow = row
		c.selectorColumn = column
	}
}

func (c *codeAreaComponent) appendCharacter(ch rune) {
	line := c.lines[c.cursorRow]
	preCursorLine := line[:c.cursorColumn]
	postCursorLine := line[c.cursorColumn:]
	c.lines[c.cursorRow] = gog.Concat(
		preCursorLine,
		[]rune{ch},
		postCursorLine,
	)
	c.cursorColumn++
}

func (c *codeAreaComponent) breakLine() {
	line := c.lines[c.cursorRow]
	preCursorLine := line[:c.cursorColumn]
	postCursorLine := line[c.cursorColumn:]
	c.lines[c.cursorRow] = preCursorLine
	c.lines = slices.Insert(c.lines, c.cursorRow+1, postCursorLine)
	c.moveCursorDown()
	c.moveCursorToStartOfLine()
}

func (c *codeAreaComponent) eraseLeft() {
	if c.cursorColumn > 0 {
		line := c.lines[c.cursorRow]
		line = slices.Delete(line, c.cursorColumn-1, c.cursorColumn)
		c.lines[c.cursorRow] = line
		c.cursorColumn--
	} else {
		if c.cursorRow > 0 {
			movedRow := c.cursorRow
			c.moveCursorUp()
			c.moveCursorToEndOfLine()
			c.lines[movedRow-1] = append(c.lines[movedRow-1], c.lines[movedRow]...)
			c.lines = slices.Delete(c.lines, movedRow, movedRow+1)
		}
	}
}

func (c *codeAreaComponent) eraseRight() {
	if c.cursorColumn < len(c.lines[c.cursorRow]) {
		line := c.lines[c.cursorRow]
		line = slices.Delete(line, c.cursorColumn, c.cursorColumn+1)
		c.lines[c.cursorRow] = line
	} else {
		if c.cursorRow < len(c.lines)-1 {
			movedRow := c.cursorRow + 1
			c.lines[movedRow-1] = append(c.lines[movedRow-1], c.lines[movedRow]...)
			c.lines = slices.Delete(c.lines, movedRow, movedRow+1)
		}
	}
}

func (c *codeAreaComponent) constructText() string {
	return strings.Join(gog.Map(c.lines, func(line []rune) string {
		return string(line)
	}), "\n")
}

func (c *codeAreaComponent) notifyChanged() {
	c.refreshTextSize()
	if c.onChange != nil {
		c.onChange(string(c.constructText()))
	}
}

func splitLines(text string) [][]rune {
	return gog.Map(strings.Split(text, "\n"), func(line string) []rune {
		return []rune(line)
	})
}

func lineToText(input []rune) string {
	return string(input)
}

func countDigits(number int) int {
	number = abs(number)

	result := 1
	for number > 9 {
		number /= 10
		result++
	}
	return result
}
