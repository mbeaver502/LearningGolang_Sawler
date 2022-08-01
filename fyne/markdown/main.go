package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	cfg.createMenuItems(win)

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

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", func() {})

	saveMenuItem := fyne.NewMenuItem("Save", func() {})
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))

	// File
	//	Open...
	//	Save
	//	Save As...
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	mainMenu := fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(mainMenu)
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			// some kind of error happened
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// user canceled save operation
			if write == nil {
				return
			}

			// save the edit widget's text to file
			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()

			// set the window title to Markdown - Filename
			app.CurrentFile = write.URI()
			win.SetTitle(fmt.Sprintf("%s - %s", win.Title(), app.CurrentFile.Name()))

			app.SaveMenuItem.Disabled = false
		}, win)

		saveDialog.Show()
	}
}
