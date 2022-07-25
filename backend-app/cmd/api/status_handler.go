package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	resp, err := json.Marshal(currentStatus)
	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
