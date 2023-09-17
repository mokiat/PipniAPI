package widget

import (
	"strconv"
	"strings"

	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"golang.org/x/exp/slices"
)

var CodeArea = co.Define(&codeAreaComponent{})

type CodeAreaData struct {
	ReadOnly bool
	Code     string
}

type CodeAreaCallbackData struct {
	OnChange func(string)
}

var defaultCodeAreaCallbackData = CodeAreaCallbackData{
	OnChange: func(string) {},
}

var _ ui.ElementRenderHandler = (*codeAreaComponent)(nil)
var _ ui.ElementKeyboardHandler = (*codeAreaComponent)(nil)

type codeAreaComponent struct {
	co.BaseComponent

	font     *ui.Font
	fontSize float32

	cursorRow           int
	cursorColumn        int
	cursorVirtualColumn int

	isReadOnly bool
	lines      [][]rune
	onChange   func(string)
}

func (c *codeAreaComponent) OnCreate() {
	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = 20.0

	c.cursorRow = 0
	c.cursorColumn = 0
}

func (c *codeAreaComponent) OnUpsert() {
	data := co.GetData[CodeAreaData](c.Properties())
	c.isReadOnly = data.ReadOnly
	c.lines = gog.Map(strings.Split(data.Code, "\n"), func(line string) []rune {
		return []rune(line)
	})

	callbackData := co.GetOptionalCallbackData[CodeAreaCallbackData](c.Properties(), defaultCodeAreaCallbackData)
	c.onChange = callbackData.OnChange

	numRows := len(c.lines)
	if c.cursorRow >= numRows {
		c.cursorRow = numRows - 1
	}
	if c.cursorColumn > len(c.lines[c.cursorRow]) {
		c.cursorColumn = len(c.lines[c.cursorRow])
	}
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
				Width:  int(contentSize.X + 100),
				Height: int(contentSize.Y),
			}),
		})
	})
}

func (c *codeAreaComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	// TODO: Take scrolling into consideration.
	// Use binary search to figure out the first and last lines that are visible.
	// This should optimize rendering of large texts.

	// TOOD: Determine correct size for container of line numbers based on the
	// number of rows and the digits.

	bounds := canvas.DrawBounds(element, false)
	isFocused := element.Window().IsElementFocused(element)

	// Background
	canvas.Reset()
	canvas.Rectangle(bounds.Position, bounds.Size)
	canvas.Fill(ui.Fill{
		Color: std.SurfaceColor,
	})

	// Lines indicator
	canvas.Reset()
	canvas.Rectangle(
		bounds.Position,
		sprec.NewVec2(90, bounds.Size.Y),
	)
	canvas.Fill(ui.Fill{
		Color: std.PrimaryLightColor,
	})

	// TODO: Draw text highlighting (if selected)

	// Draw text content
	linePosition := bounds.Position
	for i, line := range c.lines {
		lineHeight := c.font.TextSize("|", c.fontSize)

		// Draw line number
		numberPosition := sprec.Vec2Sum(linePosition, sprec.NewVec2(10.0, 0.0))
		canvas.Reset()
		canvas.FillText(strconv.Itoa(i+1), numberPosition, ui.Typography{
			Font:  c.font,
			Size:  c.fontSize,
			Color: std.OnSurfaceColor,
		})

		// Draw text
		textPosition := sprec.Vec2Sum(linePosition, sprec.NewVec2(100.0, 0.0))
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

	// Highlight
	if isFocused {
		canvas.Reset()
		canvas.SetStrokeColor(std.SecondaryColor)
		canvas.SetStrokeSize(1.0)
		canvas.Rectangle(bounds.Position, bounds.Size)
		canvas.Stroke()
	}
}

func (c *codeAreaComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Type {
	case ui.KeyboardEventTypeKeyDown, ui.KeyboardEventTypeRepeat:
		return c.onKeyboardPressEvent(element, event)

	case ui.KeyboardEventTypeType:
		return c.onKeyboardTypeEvent(element, event)

	default:
		return false
	}
}

func (c *codeAreaComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {

	case ui.KeyCodeEscape:
		element.Window().DiscardFocus()
		return true

	case ui.KeyCodeArrowUp:
		if c.isReadOnly {
			c.scrollUp()
		} else {
			c.moveCursorUp()
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowDown:
		if c.isReadOnly {
			c.scrollDown()
		} else {
			c.moveCursorDown()
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowLeft:
		if c.isReadOnly {
			c.scrollLeft()
		} else {
			c.moveCursorLeft()
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowRight:
		if c.isReadOnly {
			c.scrollRight()
		} else {
			c.moveCursorRight()
		}
		element.Invalidate()
		return true

	case ui.KeyCodeBackspace:
		if c.isReadOnly {
			return false
		}
		// TODO: Check if selection
		c.eraseLeft()
		c.onChange(c.constructText())
		element.Invalidate()
		return true

	case ui.KeyCodeDelete:
		if c.isReadOnly {
			return false
		}
		// TODO: Check if selection
		c.eraseRight()
		c.onChange(c.constructText())
		element.Invalidate()
		return true

	case ui.KeyCodeEnter:
		if c.isReadOnly {
			return false
		}
		c.breakLine()
		c.onChange(c.constructText())
		element.Invalidate()
		return true

	case ui.KeyCodeTab:
		if c.isReadOnly {
			return false
		}
		c.appendCharacter('\t')
		c.onChange(c.constructText())
		element.Invalidate()
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
	c.onChange(c.constructText())
	return true
}

func (c *codeAreaComponent) scrollUp() {
	// TODO
}

func (c *codeAreaComponent) scrollDown() {
	// TODO
}

func (c *codeAreaComponent) scrollLeft() {
	// TODO
}

func (c *codeAreaComponent) scrollRight() {
	// TODO
}

func (c *codeAreaComponent) trackVirtualColumn() {
	c.cursorVirtualColumn = c.cursorColumn
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

// TODO: Add built-in scrolling as well. The external one will not do due to auto-panning and the like.

// TODO: Mouse handler as well, so that selection is possible

// TODO: Make selection possible via keyboard as well (with SHIFT and ARROWS, etc)
