package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	output *widget.Label
}

var myApp App

func main() {
	a := app.New()
	w := a.NewWindow("Hello World!")

	// w.SetContent(widget.NewLabel("Hello World!"))
	output, entry, btn := myApp.makeUI()
	btn.Importance = widget.HighImportance

	// create a vertical box (like a WPF vertical StackPanel)
	//	output
	//	entry
	//	btn
	w.SetContent(container.NewVBox(output, entry, btn))
	w.Resize(fyne.Size{
		Width:  320,
		Height: 240,
	})

	w.ShowAndRun()
}

func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Button) {
	output := widget.NewLabel("Hello World!")
	app.output = output

	entry := widget.NewEntry()

	btn := widget.NewButton("Submit", func() {
		app.output.SetText(entry.Text)
	})

	return output, entry, btn
}
