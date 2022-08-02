package main

import (
	"database/sql"
	"goldwatcher/repository"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App            fyne.App
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	DB             repository.Repository
	MainWindow     fyne.Window
	PriceContainer *fyne.Container
	Toolbar        *widget.Toolbar
	ChartContainer *fyne.Container
	Holdings       [][]interface{}
	HoldingsTable  *widget.Table
	HTTPClient     *http.Client
}

func main() {
	var myApp Config

	// create a fyne app
	myApp.App = app.NewWithID("com.example.goldwatcher.preferences")

	// create loggers
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	myApp.HTTPClient = &http.Client{}

	// open a connection to database (SQLite)
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		myApp.ErrorLog.Panic(err)
	}
	defer sqlDB.Close()

	// create a database repository
	myApp.setupDB(sqlDB)

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

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	envPath := os.Getenv("DB_PATH")
	if envPath != "" {
		path = envPath
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
		app.InfoLog.Println("saving DB to:", path)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqlDB)

	err := app.DB.Migrate()
	if err != nil {
		app.ErrorLog.Panic(err)
	}
}
