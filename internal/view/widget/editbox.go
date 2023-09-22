package widget

import (
	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/std"
	"golang.org/x/exp/slices"
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

var defaultEditBoxCallbackData = EditBoxCallbackData{
	OnChange: func(string) {},
	OnSubmit: func(string) {},
}

var _ ui.ElementRenderHandler = (*editBoxComponent)(nil)
var _ ui.ElementKeyboardHandler = (*editBoxComponent)(nil)

type editBoxComponent struct {
	co.BaseComponent

	font     *ui.Font
	fontSize float32

	cursorColumn   int
	selectorColumn int

	isReadOnly bool
	line       []rune
	onChange   func(string)
	onSubmit   func(string)
}

func (c *editBoxComponent) OnCreate() {
	c.font = co.OpenFont(c.Scope(), "fonts/roboto-mono-regular.ttf")
	c.fontSize = 18.0

	c.cursorColumn = 0
}

func (c *editBoxComponent) OnUpsert() {
	data := co.GetData[EditBoxData](c.Properties())
	c.isReadOnly = data.ReadOnly
	c.line = []rune(data.Text)

	callbackData := co.GetOptionalCallbackData[EditBoxCallbackData](c.Properties(), defaultEditBoxCallbackData)
	c.onChange = callbackData.OnChange
	c.onSubmit = callbackData.OnSubmit

	c.cursorColumn = min(c.cursorColumn, len(c.line))
	c.selectorColumn = min(c.selectorColumn, len(c.line))
}

func (c *editBoxComponent) Render() co.Instance {
	contentSize := c.font.TextSize(string(c.line), c.fontSize)
	return co.New(std.Element, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			IdealSize: opt.V(ui.Size{
				Width:  int(contentSize.X + 10),
				Height: int(contentSize.Y),
			}),
		})
	})
}

func (c *editBoxComponent) OnRender(element *ui.Element, canvas *ui.Canvas) {
	// TODO: Take scrolling into consideration.
	// Use binary search to figure out the first and last lines that are visible.
	// This should optimize rendering of large texts.

	// TOOD: Determine correct size for container of line numbers based on the
	// number of rows and the digits.

	bounds := canvas.DrawBounds(element, false)
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
	canvas.SetClipRect(5, bounds.Width()-5, 2, bounds.Height()-2)

	textPosition := sprec.Vec2Sum(bounds.Position, sprec.NewVec2(5.0, 2.0))
	textSize := c.font.TextSize(string(c.line), c.fontSize)

	// Draw Selection
	if c.cursorColumn != c.selectorColumn {
		fromColumn := min(c.cursorColumn, c.selectorColumn)
		toColumn := max(c.cursorColumn, c.selectorColumn)
		preSelectionSize := c.font.TextSize(string(c.line[:fromColumn]), c.fontSize)
		selectionSize := c.font.TextSize(string(c.line[fromColumn:toColumn]), c.fontSize)

		selectionPosition := sprec.Vec2Sum(textPosition, sprec.NewVec2(preSelectionSize.X, 0.0))
		canvas.Reset()
		canvas.Rectangle(selectionPosition, selectionSize)
		canvas.Fill(ui.Fill{
			Color: std.SecondaryLightColor,
		})
	}

	// Draw text
	canvas.Reset()
	canvas.FillText(string(c.line), textPosition, ui.Typography{
		Font:  c.font,
		Size:  c.fontSize,
		Color: std.OnSurfaceColor,
	})

	// Draw cursor
	if !c.isReadOnly {
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
		canvas.SetStrokeColor(std.SecondaryColor)
	} else {
		canvas.SetStrokeColor(std.PrimaryColor)
	}
	canvas.SetStrokeSize(1.0)
	canvas.RoundRectangle(
		bounds.Position,
		bounds.Size,
		sprec.NewVec4(8, 8, 8, 8),
	)
	canvas.Stroke()
}

func (c *editBoxComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Type {
	case ui.KeyboardEventTypeKeyDown, ui.KeyboardEventTypeRepeat:
		return c.onKeyboardPressEvent(element, event)

	case ui.KeyboardEventTypeType:
		return c.onKeyboardTypeEvent(element, event)

	default:
		return false
	}
}

func (c *editBoxComponent) onKeyboardPressEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {

	case ui.KeyCodeEscape:
		element.Window().DiscardFocus()
		return true

	case ui.KeyCodeArrowUp:
		if !c.isReadOnly {
			c.moveCursorToStartOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowDown:
		if !c.isReadOnly {
			c.moveCursorToEndOfLine()
			if !event.Modifiers.Contains(ui.KeyModifierShift) {
				c.resetSelector()
			}
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowLeft:
		if c.isReadOnly {
			c.scrollLeft()
		} else {
			if event.Modifiers.Contains(ui.KeyModifierShift) {
				c.moveCursorLeft()
			} else {
				if c.cursorColumn == c.selectorColumn {
					c.moveCursorLeft()
				}
				c.resetSelector()
			}
		}
		element.Invalidate()
		return true

	case ui.KeyCodeArrowRight:
		if c.isReadOnly {
			c.scrollRight()
		} else {
			if event.Modifiers.Contains(ui.KeyModifierShift) {
				c.moveCursorRight()
			} else {
				if c.cursorColumn == c.selectorColumn {
					c.moveCursorRight()
				}
				c.resetSelector()
			}
		}
		element.Invalidate()
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
		c.onChange(string(c.line))
		element.Invalidate()
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
		c.onChange(string(c.line))
		element.Invalidate()
		return true

	case ui.KeyCodeEnter:
		if c.isReadOnly {
			return false
		}
		c.onSubmit(string(c.line))
		element.Invalidate()
		return true

	case ui.KeyCodeTab:
		if c.isReadOnly {
			return false
		}
		c.appendCharacter('\t')
		if !event.Modifiers.Contains(ui.KeyModifierShift) {
			c.resetSelector()
		}
		c.onChange(string(c.line))
		element.Invalidate()
		return true

	default:
		return false
	}
}

func (c *editBoxComponent) onKeyboardTypeEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	if c.isReadOnly {
		return false
	}
	c.appendCharacter(event.Rune)
	c.resetSelector()
	c.onChange(string(c.line))
	return true
}

func (c *editBoxComponent) scrollLeft() {
	// TODO
}

func (c *editBoxComponent) scrollRight() {
	// TODO
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

func (c *editBoxComponent) moveCursorToStartOfLine() {
	c.cursorColumn = 0
}

func (c *editBoxComponent) moveCursorToEndOfLine() {
	c.cursorColumn = len(c.line)
}

func (c *editBoxComponent) appendCharacter(ch rune) {
	preCursorLine := c.line[:c.cursorColumn]
	postCursorLine := c.line[c.cursorColumn:]
	c.line = gog.Concat(
		preCursorLine,
		[]rune{ch},
		postCursorLine,
	)
	c.cursorColumn++
}

func (c *editBoxComponent) eraseLeft() {
	if c.cursorColumn > 0 {
		c.line = slices.Delete(c.line, c.cursorColumn-1, c.cursorColumn)
		c.cursorColumn--
	}
}

func (c *editBoxComponent) eraseRight() {
	if c.cursorColumn < len(c.line) {
		c.line = slices.Delete(c.line, c.cursorColumn, c.cursorColumn+1)
	}
}

// TODO: Add built-in scrolling as well. The external one will not do due to auto-panning and the like.

// TODO: Mouse handler as well, so that selection is possible
