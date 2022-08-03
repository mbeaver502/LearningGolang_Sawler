package main

import (
	"fmt"
	"goldwatcher/repository"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) holdingsTab() *fyne.Container {
	app.Holdings = app.getHoldingsSlice()
	app.HoldingsTable = app.getHoldingsTable()

	holdingsContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, app.HoldingsTable),
	)

	return holdingsContainer
}

func (app *Config) getHoldingsTable() *widget.Table {
	t := widget.NewTable(
		// length
		func() (int, int) {
			// row, col
			return len(app.Holdings), len(app.Holdings[0])
		},
		// create
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		// update
		func(i widget.TableCellID, o fyne.CanvasObject) {
			// delete button goes in last cell of non-heading row
			if i.Col == len(app.Holdings[0])-1 && i.Row != 0 {
				w := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Delete?", "Are you sure?", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(app.Holdings[i.Row][0].(string))
							err := app.DB.DeleteHolding(int64(id))
							if err != nil {
								app.ErrorLog.Println(err)
							}
						}
						app.refreshHoldingsTable()
					}, app.MainWindow)
				})

				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}

			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(app.Holdings[i.Row][i.Col].(string)),
				}
			}
		},
	)

	// set the column widths
	colWidths := []float32{50, 200, 200, 200, 110}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}

	return t
}

func (app *Config) getHoldingsSlice() [][]any {
	var slice [][]any

	holdings, err := app.currentHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil
	}

	slice = append(slice, []any{"ID", "Amount", "Price", "Date", "Delete?"})

	for _, h := range holdings {
		var currentRow []any

		currentRow = append(currentRow,
			strconv.FormatInt(h.ID, 10),
			fmt.Sprintf("%d toz", h.Amount),
			fmt.Sprintf("$%.2f", float32(h.PurchasePrice/100)),
			h.PurchaseDate.Format("2006-JAN-02"))

		currentRow = append(currentRow,
			widget.NewButton("Delete", func() {}),
		)

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) currentHoldings() ([]repository.Holdings, error) {
	h, err := app.DB.AllHoldings()
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}

	return h, nil
}
