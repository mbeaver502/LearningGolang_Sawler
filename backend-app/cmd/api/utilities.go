package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, wrap string) error {
	wrapper := make(map[string]any)

	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
