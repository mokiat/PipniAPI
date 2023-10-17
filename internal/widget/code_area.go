package widget

import (
	"math"
	"strconv"
	"strings"
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
	codeAreaHistoryCapacity = 100
	codeAreaLinesWidth      = 100
	codeAreaLinesPadding    = 10
)

var CodeArea = co.Define(&codeAreaComponent{})

type CodeAreaData struct {
	ReadOnly bool
	Code     string
}

type CodeAreaCallbackData struct {
	OnChange func(string)
}

var _ ui.ElementResizeHandler = (*codeAreaComponent)(nil)
var _ ui.ElementRenderHandler = (*codeAreaComponent)(nil)
var _ ui.ElementKeyboardHandler = (*codeAreaComponent)(nil)
var _ ui.ElementMouseHandler = (*codeAreaComponent)(nil)
var _ ui.ElementHistoryHandler = (*codeAreaComponent)(nil)
var _ ui.ElementClipboardHandler = (*codeAreaComponent)(nil)

type codeAreaComponent struct {
	co.BaseComponent

	history *state.History

	font     *ui.Font
	fontSize float32

	cursorRow           int
	cursorColumn        int
	cursorVirtualColumn int

	selectorRow    int
	selectorColumn int

	isReadOnly bool
	lines      [][]rune
	onChange   func(string)

	textWidth  int
	textHeight int

	offsetX    int
	offsetY    int
	maxOffsetX int
	maxOffsetY int

	isDragging bool
}

func (c *codeAreaComponent) OnCreate() {
	c.history = state.NewHistory(editBoxHistoryCapacity)

	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = 18.0 // TODO: Was 20. Is it better with 18? EditBox is 18.

	c.cursorRow = 0
	c.cursorColumn = 0
	c.cursorVirtualColumn = 0
}

func (c *codeAreaComponent) OnUpsert() {
	data := co.GetData[CodeAreaData](c.Properties())
	c.isReadOnly = data.ReadOnly
	if data.Code != c.constructText() {
		c.history.Clear()
	}
	c.lines = gog.Map(strings.Split(data.Code, "\n"), func(line string) []rune {
		return []rune(line)
	})
	c.refreshTextSize()

	callbackData := co.GetOptionalCallbackData[CodeAreaCallbackData](c.Properties(), CodeAreaCallbackData{})
	c.onChange = callbackData.OnChange

	c.cursorRow = min(c.cursorRow, len(c.lines)-1)
	c.cursorColumn = min(c.cursorColumn, len(c.lines[c.cursorRow]))
	c.selectorRow = min(c.selectorRow, len(c.lines)-1)
	c.selectorColumn = min(c.selectorColumn, len(c.lines[c.selectorRow]))
}

func (c *codeAreaComponent) Render() co.Instance {
	var contentSize sprec.Vec2
	for _, line := range c.lines {
		textSize := c.font.TextSize(string(line), c.fontSize)
		contentSize.X = max(contentSize.X, textSize.X)
		contentSize.Y += textSize.Y
	}
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			IdealSize: opt.V(ui.Size{
				Width:  int(codeAreaLinesWidth + codeAreaLinesPadding + contentSize.X),
				Height: int(contentSize.Y),
			}),
		})
	})
}

func (c *codeAreaComponent) refreshTextSize() {
	var textSize sprec.Vec2
	for _, line := range c.lines {
		lineSize := c.font.TextSize(string(line), c.fontSize)
		textSize.X = max(textSize.X, lineSize.X)
		textSize.Y += lineSize.Y
	}
	c.textWidth = int(math.Ceil(float64(textSize.X)))
	c.textHeight = int(math.Ceil(float64(textSize.Y)))
}

func (c *codeAreaComponent) refreshScrollBounds(element *ui.Element) {
	bounds := element.ContentBounds()
	availableTextWidth := bounds.Width - 110
	availableTextHeight := bounds.Height
	c.maxOffsetX = max(c.textWidth-availableTextWidth, 0)
	c.maxOffsetY = max(c.textHeight-availableTextHeight, 0)
	c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
	c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
}

func (c *codeAreaComponent) OnResize(element *ui.Element, bounds ui.Bounds) {
	c.refreshScrollBounds(element)
}

func (c *codeAreaComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	c.refreshScrollBounds(element)

	// TODO: Take scrolling into consideration.
	// Use binary search to figure out the first and last lines that are visible.
	// This should optimize rendering of large texts.

	// TODO: Determine correct size for container of line numbers based on the
	// number of rows and the digits.

	bounds := canvas.DrawBounds(element, false)
	isFocused := element.Window().IsElementFocused(element)

	// Background
	canvas.Reset()
	canvas.Rectangle(bounds.Position, bounds.Size)
	canvas.Fill(ui.Fill{
		Color: std.SurfaceColor,
	})

	selection := c.selection()

	// Draw text content
	lineHeight := c.font.TextSize("|", c.fontSize)
	linePosition := sprec.Vec2Diff(bounds.Position, sprec.NewVec2(0.0, float32(c.offsetY)))

	for i, line := range c.lines {
		textPosition := sprec.Vec2Sum(linePosition, sprec.NewVec2(100.0-float32(c.offsetX), 0.0))

		// Draw Selection
		if selection.ContainsRow(i) {
			fromColumn, toColumn := selection.ColumnSpan(i, len(line))
			preSelectionSize := c.font.TextSize(string(line[:fromColumn]), c.fontSize)
			selectionSize := c.font.TextSize(string(line[fromColumn:toColumn]), c.fontSize)

			selectionPosition := sprec.Vec2Sum(textPosition, sprec.NewVec2(preSelectionSize.X, 0.0))
			canvas.Reset()
			canvas.Rectangle(selectionPosition, selectionSize)
			canvas.Fill(ui.Fill{
				Color: std.SecondaryLightColor,
			})
		}

		// Draw text
		canvas.Reset()
		canvas.FillText(string(line), textPosition, ui.Typography{
			Font:  c.font,
			Size:  c.fontSize,
			Color: std.OnSurfaceColor,
		})

		// Draw cursor
		if i == c.cursorRow && !c.isReadOnly {
			cursorColumn := min(c.cursorColumn, len(line))
			preCursorText := line[:cursorColumn]
			preCursorTextSize := c.font.TextSize(string(preCursorText), c.fontSize)

			// TODO: Take tilt into account and use like stroke instead of rect fill.
			cursorPosition := sprec.Vec2Sum(textPosition, sprec.NewVec2(preCursorTextSize.X, 0.0))
			cursorWidth := float32(1.0)
			canvas.Reset()
			canvas.Rectangle(cursorPosition, sprec.NewVec2(cursorWidth, lineHeight.Y))
			canvas.Fill(ui.Fill{
				Color: std.PrimaryColor,
			})
		}

		linePosition.Y += lineHeight.Y
	}

	// Lines indicator
	canvas.Reset()
	canvas.Rectangle(
		bounds.Position,
		sprec.NewVec2(90, bounds.Size.Y),
	)
	canvas.Fill(ui.Fill{
		Color: std.PrimaryLightColor,
	})

	linePosition = sprec.Vec2Diff(bounds.Position, sprec.NewVec2(0.0, float32(c.offsetY)))
	for i := range c.lines {
		// Draw line number
		numberPosition := sprec.Vec2Sum(linePosition, sprec.NewVec2(10.0, 0.0))
		canvas.Reset()
		canvas.FillText(strconv.Itoa(i+1), numberPosition, ui.Typography{
			Font:  c.font,
			Size:  c.fontSize,
			Color: std.OnSurfaceColor,
		})
		linePosition.Y += lineHeight.Y
	}

	// Highlight
	if isFocused {
		canvas.Reset()
		canvas.SetStrokeColor(std.SecondaryColor)
		canvas.SetStrokeSize(1.0)
		canvas.Rectangle(bounds.Position, bounds.Size)
		canvas.Stroke()
	}
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
	// case ui.ClipboardActionCut:
	// 	if c.isReadOnly {
	// 		return false
	// 	}
	// 	if c.hasSelection() {
	// 		text := string(c.selectedText())
	// 		element.Window().RequestCopy(text)
	// 		c.history.Do(c.changeDeleteSelection())
	// 		c.notifyChanged()
	// 	}
	// 	return true

	// case ui.ClipboardActionCopy:
	// 	if c.hasSelection() {
	// 		text := string(c.selectedText())
	// 		element.Window().RequestCopy(text)
	// 	}
	// 	return true

	// case ui.ClipboardActionPaste:
	// 	if c.isReadOnly {
	// 		return false
	// 	}
	// 	if c.hasSelection() {
	// 		c.history.Do(c.changeReplaceSelection([]rune(event.Text)))
	// 	} else {
	// 		c.history.Do(c.changeAppendText([]rune(event.Text)))
	// 	}
	// 	c.notifyChanged()
	// 	return true

	default:
		return false
	}
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
			c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
		} else {
			c.offsetX -= event.ScrollX
			c.offsetX = min(max(c.offsetX, 0), c.maxOffsetX)
			c.offsetY -= event.ScrollY
			c.offsetY = min(max(c.offsetY, 0), c.maxOffsetY)
		}
		element.Invalidate()
		return true

	// case ui.MouseActionDown:
	// 	if event.Button != ui.MouseButtonLeft {
	// 		return false
	// 	}
	// 	c.isDragging = true
	// 	c.cursorColumn = c.findCursorColumn(element, event.X)
	// 	if !event.Modifiers.Contains(ui.KeyModifierShift) {
	// 		c.resetSelector()
	// 	}
	// 	element.Invalidate()
	// 	return true

	// case ui.MouseActionMove: // TODO: Use dragging event
	// 	if c.isDragging {
	// 		c.cursorColumn = c.findCursorColumn(element, event.X)
	// 		element.Invalidate()
	// 	}
	// 	return true

	// case ui.MouseActionUp:
	// 	if event.Button != ui.MouseButtonLeft {
	// 		return false
	// 	}
	// 	if c.isDragging {
	// 		c.isDragging = false
	// 		c.cursorColumn = c.findCursorColumn(element, event.X)
	// 		element.Invalidate()
	// 	}
	// return true

	default:
		return false
	}
}

func (c *codeAreaComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
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
		element.Window().DiscardFocus()
		return true

	case ui.KeyCodeArrowUp:
		if c.isReadOnly {
			c.scrollUp()
		} else {
			c.moveCursorUp()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		return true

	case ui.KeyCodeArrowDown:
		if c.isReadOnly {
			c.scrollDown()
		} else {
			c.moveCursorDown()
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
				if !c.selection().Valid() {
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
				if !c.selection().Valid() {
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
		// TODO: Check if selection
		c.eraseLeft()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		c.onChange(c.constructText())
		return true

	case ui.KeyCodeDelete:
		if c.isReadOnly {
			return false
		}
		// TODO: Check if selection
		c.eraseRight()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		c.onChange(c.constructText())
		return true

	case ui.KeyCodeEnter:
		if c.isReadOnly {
			return false
		}
		c.breakLine()
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		c.onChange(c.constructText())
		return true

	case ui.KeyCodeTab:
		if c.isReadOnly {
			return false
		}
		c.appendCharacter('\t')
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		c.onChange(c.constructText())
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
	c.resetSelector()
	c.onChange(c.constructText())
	return true
}

func (c *codeAreaComponent) hasSelection() bool {
	return c.cursorColumn != c.selectorColumn ||
		c.cursorRow != c.selectorRow
}

func (c *codeAreaComponent) selectAll() {
	c.selectorRow = 0
	c.selectorColumn = 0
	c.cursorRow = len(c.lines) - 1
	c.cursorColumn = len(c.lines[c.cursorRow])
	c.cursorVirtualColumn = c.cursorColumn
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

func (c *codeAreaComponent) trackVirtualColumn() {
	c.cursorVirtualColumn = c.cursorColumn
}

func (c *codeAreaComponent) resetSelector() {
	c.selectorRow = c.cursorRow
	c.selectorColumn = c.cursorColumn
}

func (c *codeAreaComponent) moveCursorUp() {
	if c.cursorRow > 0 {
		c.cursorRow--
		c.cursorColumn = c.cursorVirtualColumn
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
		c.cursorColumn = c.cursorVirtualColumn
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
	c.trackVirtualColumn()
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
	c.trackVirtualColumn()
}

func (c *codeAreaComponent) moveCursorToStartOfLine() {
	c.cursorColumn = 0
	c.trackVirtualColumn()
}

func (c *codeAreaComponent) moveCursorToEndOfLine() {
	c.cursorColumn = len(c.lines[c.cursorRow])
	c.trackVirtualColumn()
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
	c.trackVirtualColumn()
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
		c.trackVirtualColumn()
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

func (c *codeAreaComponent) selection() selectionSpan {
	switch {
	case c.cursorRow < c.selectorRow:
		return selectionSpan{
			FromRow:    c.cursorRow,
			ToRow:      c.selectorRow,
			FromColumn: c.cursorColumn,
			ToColumn:   c.selectorColumn,
		}
	case c.selectorRow < c.cursorRow:
		return selectionSpan{
			FromRow:    c.selectorRow,
			ToRow:      c.cursorRow,
			FromColumn: c.selectorColumn,
			ToColumn:   c.cursorColumn,
		}
	default:
		return selectionSpan{
			FromRow:    c.cursorRow,
			ToRow:      c.cursorRow,
			FromColumn: min(c.cursorColumn, c.selectorColumn),
			ToColumn:   max(c.cursorColumn, c.selectorColumn),
		}
	}
}

func (c *codeAreaComponent) notifyChanged() {
	if c.onChange != nil {
		c.refreshTextSize()
		c.onChange(string(c.constructText()))
	}
}

type selectionSpan struct {
	FromRow    int
	FromColumn int
	ToRow      int
	ToColumn   int
}

func (s selectionSpan) Valid() bool {
	return s.FromRow != s.ToRow || s.FromColumn != s.ToColumn
}

func (s selectionSpan) ContainsRow(row int) bool {
	return s.FromRow <= row && row <= s.ToRow
}

func (s selectionSpan) ColumnSpan(row, lineLength int) (int, int) {
	if row == s.FromRow && row == s.ToRow {
		return s.FromColumn, s.ToColumn
	}
	switch row {
	case s.FromRow:
		return s.FromColumn, lineLength
	case s.ToRow:
		return 0, s.ToColumn
	default:
		return 0, lineLength
	}
}

type codeAreaChange struct {
	when    time.Time
	forward []func()
	reverse []func()
}

func (c *codeAreaChange) Apply() {
	for _, action := range c.forward {
		action()
	}
}

func (c *codeAreaChange) Revert() {
	for _, action := range c.reverse {
		action()
	}
}

func (c *codeAreaChange) Extend(other state.Change) bool {
	otherChange, ok := other.(*codeAreaChange)
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
