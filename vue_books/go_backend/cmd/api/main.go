package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) serve() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Ok      bool   `json:"ok"`
			Message string `json:"message"`
		}

		payload.Ok = true
		payload.Message = "hello world!"

		out, err := json.Marshal(payload)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	})

	app.infoLog.Println("API listening on port", app.config.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), nil)
}
