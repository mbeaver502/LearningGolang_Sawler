package main

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (app *Config) makeUI() {
	// get current price of gold
	openPrice, currentPrice, priceChange := app.getPriceText()

	// put price information into container
	priceContent := container.NewGridWithColumns(3, openPrice, currentPrice, priceChange)
	app.PriceContainer = priceContent

	// get toolbar
	toolbar := app.getToolbar()
	app.Toolbar = toolbar

	// get tabs
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Prices", theme.HomeIcon(), canvas.NewText("price content goes here...", nil)),
		container.NewTabItemWithIcon("Holdings", theme.InfoIcon(), canvas.NewText("holdings content goes here...", nil)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// add container to window
	finalContent := container.NewVBox(
		priceContent,
		toolbar,
		tabs,
	)

	app.MainWindow.SetContent(finalContent)
}
