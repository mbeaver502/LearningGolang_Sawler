package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON reads JSON from the request into the variable pointed to by data.
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	// make sure our request had only *one* JSON object in it
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must have at most one JSON value")
	}

	return nil
}

// writeJSON writes the given data to the response with the supplied status code and optional header.
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		// use only the first supplied header
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON writes the given error to the response with the optional status code.
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJSON(w, statusCode, payload)
}
