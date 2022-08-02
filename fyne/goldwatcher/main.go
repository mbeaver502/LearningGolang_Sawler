package main

import (
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Config struct {
	App            fyne.App
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	MainWindow     fyne.Window
	PriceContainer *fyne.Container
	HTTPClient     *http.Client
}

var myApp Config

func main() {
	// create a fyne app
	myApp.App = app.NewWithID("com.example.goldwatcher.preferences")

	// create loggers
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	myApp.HTTPClient = &http.Client{}

	// open a connection to database (SQLite)
	// create a database repository

	// create and size a fyne window
	myApp.MainWindow = myApp.App.NewWindow("GoldWatcher")
	myApp.MainWindow.Resize(fyne.NewSize(770, 410))

	myApp.makeUI()

	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()
	myApp.MainWindow.CenterOnScreen()

	// show and run app
	myApp.MainWindow.ShowAndRun()
}
