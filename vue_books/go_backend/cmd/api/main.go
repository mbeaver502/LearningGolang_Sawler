package main

import (
	"books_backend/internal/data"
	"books_backend/internal/driver"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
}

type application struct {
	config   config
	models   data.Models
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		models:   data.New(db.SQL),
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ErrorLog:          app.errorLog,
	}

	return srv.ListenAndServe()
}
