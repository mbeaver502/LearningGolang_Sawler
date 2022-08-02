package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (app *Config) pricesTab() *fyne.Container {
	chart := app.getChart()
	chartContainer := container.NewVBox(chart)
	app.ChartContainer = chartContainer

	return chartContainer
}

func (app *Config) getChart() *canvas.Image {
	apiURL := fmt.Sprintf("https://goldprice.org/charts/gold_3d_b_o_%s_x.png", strings.ToLower(currency))

	var img *canvas.Image

	err := app.downloadFile(apiURL, "gold.png")
	if err != nil {
		// use bundled image
		img = canvas.NewImageFromResource(resourceUnreachablePng)
	} else {
		img = canvas.NewImageFromFile("gold.png")
	}

	img.SetMinSize(fyne.Size{
		Width:  app.MainWindow.Canvas().Size().Width,
		Height: app.MainWindow.Canvas().Size().Height,
	})

	// fill available space while keeping original aspect
	img.FillMode = canvas.ImageFillOriginal

	return img
}

func (app *Config) downloadFile(URL string, filename string) error {
	// get the response bytes from calling the url
	resp, err := app.HTTPClient.Get(URL)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	// read the response body into a slice of bytes
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// decode the bytes
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	// create a file on the local filesystem
	out, err := os.Create(fmt.Sprintf("./%s", filename))
	if err != nil {
		return err
	}

	// encode the decoded bytes into a PNG in the created file
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
