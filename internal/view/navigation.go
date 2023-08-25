package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func (w *Window) newNavigationPanel() fyne.CanvasObject {
	top := w.newEnvironmentSelection()
	middle := w.newResourceTree()

	return container.NewBorder(top, nil, nil, nil, middle)
}
