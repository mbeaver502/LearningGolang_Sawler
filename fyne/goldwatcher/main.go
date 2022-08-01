package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Config struct {
	App      fyne.App
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var myApp Config

func main() {
	// create a fyne app
	myApp.App = app.NewWithID("com.example.goldwatcher.preferences")

	// create loggers
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a connection to database (SQLite)
	// create a database repository

	// create and size a fyne window
	win := myApp.App.NewWindow("GoldWatcher")

	// show and run app
	win.ShowAndRun()
}
