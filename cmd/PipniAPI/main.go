package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/mokiat/PipniAPI/internal/fyneutil"
	"github.com/mokiat/PipniAPI/internal/mvc"
	"github.com/mokiat/PipniAPI/internal/view"
	"github.com/mokiat/PipniAPI/resources"
)

func main() {
	if err := runApp(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runApp() error {
	icon, err := fyneutil.LoadResourceFromFS(resources.FS, "images/icon.png")
	if err != nil {
		return fmt.Errorf("error loading icon: %w", err)
	}

	eventBus := mvc.NewEventBus()

	a := app.NewWithID("com.mokiat.pipniapi")
	w := a.NewWindow("PipniAPI")
	w.SetIcon(icon)
	w.SetMaster()
	w.SetContent(view.NewWindow(eventBus))
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
	return nil
}
