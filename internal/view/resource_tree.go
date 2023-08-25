package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (w *Window) newResourceTree() fyne.CanvasObject {
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

	return apiTree
}
