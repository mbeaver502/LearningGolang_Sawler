package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func (app *Config) getPriceText() (*canvas.Text, *canvas.Text, *canvas.Text) {
	var g Gold
	var open, current, change *canvas.Text

	g.Client = app.HTTPClient

	p, err := g.GetPrices()
	if err != nil {
		app.ErrorLog.Println(err)

		gray := color.NRGBA{
			R: 155,
			G: 155,
			B: 155,
			A: 255,
		}

		open = canvas.NewText("Open: Unreachable", gray)
		current = canvas.NewText("Current: Unreachable", gray)
		change = canvas.NewText("Change: Unreachable", gray)

		return open, current, change
	}

	// default display color
	displayColor := color.NRGBA{
		R: 0,
		G: 180,
		B: 0,
		A: 255,
	}

	if p.Price < p.PreviousClose {
		displayColor = color.NRGBA{
			R: 180,
			G: 0,
			B: 0,
			A: 255,
		}
	}

	openText := fmt.Sprintf("Open: $%.4f %s", p.PreviousClose, currency)
	currentText := fmt.Sprintf("Current: $%.4f %s", p.Price, currency)
	changeText := fmt.Sprintf("Change: $%.4f %s", p.ChangeInPrice, currency)

	open = canvas.NewText(openText, nil)
	current = canvas.NewText(currentText, displayColor)
	change = canvas.NewText(changeText, displayColor)

	open.Alignment = fyne.TextAlignLeading
	current.Alignment = fyne.TextAlignCenter
	change.Alignment = fyne.TextAlignTrailing

	return open, current, change
}
