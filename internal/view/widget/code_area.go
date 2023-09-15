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

	rowCount    int
	columnCount int

	cursorRow    int
	cursorColumn int

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

	c.rowCount = len(c.lines)
	c.columnCount = 0
	for _, line := range c.lines {
		c.columnCount = max(c.columnCount, len(line))
	}

	if c.cursorRow >= c.rowCount {
		c.cursorRow = c.rowCount - 1
	}
	if c.cursorColumn > c.columnCount {
		c.cursorColumn = c.columnCount
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
	if c.isReadOnly {
		return false
	}

	switch event.Type {
	case ui.KeyboardEventTypeKeyDown, ui.KeyboardEventTypeRepeat:
		switch event.Code {
		case ui.KeyCodeEscape:
			element.Window().DiscardFocus()
		case ui.KeyCodeArrowUp:
			if c.moveCursorUp() {
				element.Invalidate()
				return true
			}
		case ui.KeyCodeArrowDown:
			if c.moveCursorDown() {
				element.Invalidate()
				return true
			}
		case ui.KeyCodeArrowLeft:
			if c.moveCursorLeft() {
				element.Invalidate()
				return true
			}
		case ui.KeyCodeArrowRight:
			if c.moveCursorRight() {
				element.Invalidate()
				return true
			}
		case ui.KeyCodeBackspace:
			if c.eraseLeft() {
				c.onChange(c.constructText())
				c.Invalidate()
				element.Invalidate()
				return true
			}
		case ui.KeyCodeDelete:
			if c.eraseRight() {
				c.onChange(c.constructText())
				c.Invalidate()
				element.Invalidate()
				return true
			}
		case ui.KeyCodeEnter:
			if c.splitLine() {
				c.onChange(c.constructText())
				c.Invalidate()
				element.Invalidate()
				return true
			}
		}

	case ui.KeyboardEventTypeType:
		if c.appendCharacter(event.Rune) {
			c.onChange(c.constructText())
			c.Invalidate()
			element.Invalidate()
			return true
		}

	}

	return false
}

func (c *codeAreaComponent) moveCursorUp() bool {
	if c.cursorRow == 0 {
		return false
	}
	c.cursorRow--
	return true
}

func (c *codeAreaComponent) moveCursorDown() bool {
	if c.cursorRow >= c.rowCount-1 {
		return false
	}
	c.cursorRow++
	return true
}

func (c *codeAreaComponent) moveCursorLeft() bool {
	if c.cursorColumn == 0 {
		return false
	}
	c.cursorColumn--
	return true
}

func (c *codeAreaComponent) moveCursorRight() bool {
	if c.cursorColumn >= c.columnCount {
		return false
	}
	c.cursorColumn++
	return true
}

func (c *codeAreaComponent) appendCharacter(ch rune) bool {
	line := c.lines[c.cursorRow]
	cursorColumn := min(c.cursorColumn, len(line))
	preCursorLine := line[:cursorColumn]
	postCursorLine := line[cursorColumn:]
	line = gog.Concat(
		preCursorLine,
		[]rune{ch},
		postCursorLine,
	)
	c.lines[c.cursorRow] = line

	c.columnCount = max(c.columnCount, len(line))
	c.cursorColumn++
	return true
}

func (c *codeAreaComponent) splitLine() bool {
	line := c.lines[c.cursorRow]
	cursorColumn := min(c.cursorColumn, len(line))
	preCursorLine := line[:cursorColumn]
	postCursorLine := line[cursorColumn:]
	c.lines[c.cursorRow] = preCursorLine
	c.lines = slices.Insert(c.lines, c.cursorRow+1, postCursorLine)
	c.rowCount++

	c.cursorRow++
	c.cursorColumn = 0
	return true
}

func (c *codeAreaComponent) eraseLeft() bool {
	return false // TODO
}

func (c *codeAreaComponent) eraseRight() bool {
	return false // TODO
}

func (c *codeAreaComponent) constructText() string {
	return strings.Join(gog.Map(c.lines, func(line []rune) string {
		return string(line)
	}), "\n")
}

// TODO: Add built-in scrolling as well. The external one will not do due to auto-panning and the like.

// TODO: Mouse handler as well, so that selection is possible

// TODO: Make selection possible via keyboard as well (with SHIFT and ARROWS, etc)
