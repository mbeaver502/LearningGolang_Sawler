package main

import (
	"net/http"
)

// Broker is the default handler.
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "broker says hello world!",
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
