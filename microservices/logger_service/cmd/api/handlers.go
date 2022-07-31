package main

import (
	"log"
	"loggerservice/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog writes a log entry from the request.
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read JSON request
	var requestPayload jsonPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// insert the data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// write back a response
	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
