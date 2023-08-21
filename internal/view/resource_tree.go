package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func NewResourceTree(eventBus *mvc.EventBus, mdl *model.Registry) fyne.CanvasObject {
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
