package main

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) getToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		// horizontal spacer, starting at left
		widget.NewToolbarSpacer(),
		// pencil button
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}),
		// refresh button
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {}),
		// settings cog button
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
	)

	return toolbar
}
