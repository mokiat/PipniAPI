package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewWindow(eventBus *mvc.EventBus) fyne.CanvasObject {
	return container.NewHSplit(
		NewTreeMenu(eventBus),
		NewContentArena(eventBus),
	)
}

func NewTreeMenu(eventBus *mvc.EventBus) fyne.CanvasObject {
	envSelector := widget.NewSelect([]string{"Staging", "Production"}, func(string) {
		// TOOD
	})
	envSelector.SetSelected("Staging")

	settings := widget.NewButton("Settings", func() {})

	top := container.NewBorder(nil, nil, nil, settings, envSelector)

	apiTree := widget.NewTree(
		func(parentID widget.TreeNodeID) []widget.TreeNodeID {
			if parentID == "" {
				return []widget.TreeNodeID{"User", "Admin"}
			}
			return nil
		},
		func(id widget.TreeNodeID) bool {
			return id == ""
		},
		func(bool) fyne.CanvasObject {
			return widget.NewLabel("Resource Tree")
		},
		func(id widget.TreeNodeID, isBranch bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(id)
		},
	)
	apiTree.Select("User")

	// apiTree = widget.NewTreeWithStrings(map[string][]string{
	// 	"":      {"hello", "world"},
	// 	"hello": {"1", "2"},
	// })

	bottom := widget.NewCheck("Hello", func(bool) {})

	return container.NewBorder(top, bottom, nil, nil, apiTree)
}

func NewContentArena(eventBus *mvc.EventBus) fyne.CanvasObject {
	return container.NewVSplit(
		widget.NewLabel("request"),
		widget.NewLabel("response"),
	)
}
