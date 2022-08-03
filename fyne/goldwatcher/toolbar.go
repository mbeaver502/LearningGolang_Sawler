package main

import (
	"goldwatcher/repository"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) getToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		// horizontal spacer, starting at left
		widget.NewToolbarSpacer(),
		// pencil button
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			app.addHoldingsDialog()
		}),
		// refresh button
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			app.refreshPriceContent()
		}),
		// settings cog button
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
	)

	return toolbar
}

func (app *Config) addHoldingsDialog() dialog.Dialog {
	addAmountEntry := widget.NewEntry()
	purchaseDateEntry := widget.NewEntry()
	purchasePriceEntry := widget.NewEntry()

	app.AddHoldingsPurchaseAmountEntry = addAmountEntry
	app.AddHoldingsPurchaseDateEntry = purchaseDateEntry
	app.AddHoldingsPurchasePriceEntry = purchasePriceEntry

	purchaseDateEntry.Validator = func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}

	addAmountEntry.Validator = func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return err
		}
		return nil
	}

	purchasePriceEntry.Validator = func(s string) error {
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return err
		}
		return nil
	}

	purchaseDateEntry.SetPlaceHolder("YYYY-MM-DD")

	d := dialog.NewForm(
		"Add Holding",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{
				Text:   "Amount in toz",
				Widget: addAmountEntry,
			},
			{
				Text:   "Purchase Price",
				Widget: purchasePriceEntry,
			},
			{
				Text:   "Purchase Date",
				Widget: purchaseDateEntry,
			},
		},
		func(valid bool) {
			if valid {
				amount, _ := strconv.Atoi(addAmountEntry.Text)
				purchaseDate, _ := time.Parse("2006-01-02", purchaseDateEntry.Text)
				purchasePrice, _ := strconv.ParseFloat(purchasePriceEntry.Text, 64)
				purchasePrice *= 100.0

				_, err := app.DB.InsertHolding(repository.Holdings{
					Amount:        amount,
					PurchaseDate:  purchaseDate,
					PurchasePrice: int(purchasePrice),
				})
				if err != nil {
					app.ErrorLog.Println(err)
				}

				app.refreshHoldingsTable()
			}
		},
		app.MainWindow,
	)

	d.Resize(fyne.Size{Width: 400})
	d.Show()

	return d
}
