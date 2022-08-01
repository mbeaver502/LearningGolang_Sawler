package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// what we expect to receive in request
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// try to parse the request body
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate that the user actually exists
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// validate the user's password
	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if !valid || err != nil {
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log the request using Logger service
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name string, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
		return err
	}

	// this is the name of the Logger service running inside Docker
	// see the Docker-Compose file
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}

	// attempt to call Logger service with request
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
