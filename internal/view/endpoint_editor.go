package view

import (
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mokiat/PipniAPI/internal/model"
	"github.com/mokiat/PipniAPI/internal/mvc"
)

func (w *Window) newEndpointEditor(mdlEditor *model.EndpointEditor) fyne.CanvasObject {
	methodSelectControl := widget.NewSelect(nil, nil)

	methodSelectControl.Options = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	updateSelectedMethod := func(method string) {
		methodSelectControl.SetSelected(method)
	}
	updateSelectedMethod(mdlEditor.Method())

	methodSelectControl.OnChanged = func(value string) {
		mdlEditor.SetMethod(value)
	}

	goButton := widget.NewButton("GO", nil)
	goButton.OnTapped = func() {
		log.Printf("MAKING REQUEST [%s] %s", mdlEditor.Method(), mdlEditor.URI())
	}

	uriInput := widget.NewEntry()

	updateURIInput := func(uri string) {
		uriInput.SetText(uri)
	}
	updateURIInput(mdlEditor.URI())

	uriInput.OnChanged = func(value string) {
		mdlEditor.SetURI(value)
	}

	w.eventBus.Subscribe(func(event mvc.Event) {
		// TODO: Handle editor close by unsubscribing and deleting widgets
		switch event := event.(type) {
		case model.EndpointMethodChangedEvent:
			if event.Editor == mdlEditor {
				updateSelectedMethod(event.Method)
			}
		case model.EndpointURIChangedEvent:
			if event.Editor == mdlEditor {
				updateURIInput(event.URI)
			}
		}
	})

	top := container.NewBorder(nil, nil, methodSelectControl, goButton, uriInput)

	requestBody := widget.NewMultiLineEntry()
	requestBody.SetText("Hello World")

	responseBody := widget.NewRichText()
	responseBody.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Style: widget.RichTextStyleInline,
			Text:  "Response...",
		},
	}
	responseBody.Refresh()

	content := container.NewHSplit(
		requestBody,
		responseBody,
	)
	return container.NewBorder(top, nil, nil, nil, content)
}
