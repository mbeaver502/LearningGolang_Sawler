package main

import (
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
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// validate the user's password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if !valid || err != nil {
		log.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
