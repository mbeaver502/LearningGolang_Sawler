package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()

	// create a window on the app
	win := a.NewWindow("Markdown")

	// get the user interface for the window
	edit, preview := cfg.makeUI()

	// set the contents of the window
	// editor on left, preview on right
	win.SetContent(container.NewHSplit(edit, preview))

	// show the window and run the app
	win.Resize(fyne.Size{
		Width:  640,
		Height: 480,
	})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	// create an event listener such that
	// when the editor's text changes
	// the markdown is rendered in preview
	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}
