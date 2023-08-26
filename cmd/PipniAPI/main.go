package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/mokiat/PipniAPI/internal/fyneutil"
	"github.com/mokiat/PipniAPI/internal/view"
	"github.com/mokiat/PipniAPI/resources"
)

func main() {
	icon := fyneutil.Must(fyneutil.LoadResourceFromFS(resources.FS, "images/icon.png"))

	a := app.NewWithID("com.mokiat.pipniapi")
	a.SetIcon(icon)

	w := a.NewWindow("PipniAPI")
	pipniWindow := view.NewWindow(a, w)
	w.SetMainMenu(pipniWindow.RenderMainMenu())
	w.SetContent(pipniWindow.RenderContent())
	w.SetIcon(icon)
	w.SetMaster()
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}
