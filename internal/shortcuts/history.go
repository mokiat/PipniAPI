package shortcuts

import (
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/ui"
)

func IsUndo(os app.OS, event ui.KeyboardEvent) bool {
	switch os {
	case app.OSDarwin:
		return IsKeyCombo(event,
			ui.KeyModifiers(ui.KeyModifierSuper), ui.KeyCodeZ,
		)
	case app.OSWindows:
		fallthrough
	case app.OSLinux:
		fallthrough
	default:
		return IsKeyCombo(event,
			ui.KeyModifiers(ui.KeyModifierControl), ui.KeyCodeZ,
		)
	}
}

func IsRedo(os app.OS, event ui.KeyboardEvent) bool {
	switch os {
	case app.OSDarwin:
		return IsKeyCombo(event,
			ui.KeyModifiers(ui.KeyModifierSuper, ui.KeyModifierShift), ui.KeyCodeZ,
		)
	case app.OSWindows:
		return IsKeyCombo(event,
			ui.KeyModifiers(ui.KeyModifierControl), ui.KeyCodeY,
		)
	case app.OSLinux:
		fallthrough
	default:
		return IsKeyCombo(event,
			ui.KeyModifiers(ui.KeyModifierControl, ui.KeyModifierShift), ui.KeyCodeZ,
		)
	}
}

func IsKeyCombo(event ui.KeyboardEvent, modifiers ui.KeyModifierSet, code ui.KeyCode) bool {
	return event.Modifiers == modifiers && event.Code == code
}
