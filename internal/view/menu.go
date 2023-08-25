package view

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (w *Window) newMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		w.newFileMenu(),
		w.newEditMenu(),
		w.newHelpMenu(),
	)
}

func (w *Window) newFileMenu() *fyne.Menu {
	save := fyne.NewMenuItem("Save", func() {

	})

	return fyne.NewMenu("File",
		save,
	)
}

func (w *Window) newEditMenu() *fyne.Menu {
	undo := fyne.NewMenuItem("Undo", func() {

	})

	redo := fyne.NewMenuItem("Redo", func() {

	})

	cut := fyne.NewMenuItem("Cut", func() {

	})

	copy := fyne.NewMenuItem("Copy", func() {

	})

	paste := fyne.NewMenuItem("Paste", func() {

	})

	return fyne.NewMenu("Edit",
		undo,
		redo,
		fyne.NewMenuItemSeparator(),
		cut,
		copy,
		paste,
	)
}

func (w *Window) newHelpMenu() *fyne.Menu {
	github := fyne.NewMenuItem("Source Code", func() {
		url, _ := url.Parse("https://github.com/mokiat/PipniAPI")
		w.app.OpenURL(url)
	})

	about := fyne.NewMenuItem("About", func() {
		dialog.ShowCustom("About", "OK", widget.NewLabel("An open-source tool for making API calls."), w.win)
	})

	return fyne.NewMenu("Help",
		github,
		fyne.NewMenuItemSeparator(),
		about,
	)
}
