package main

import (
	"errors"
	"net/http"
	"time"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

// Login handles logging in a user.
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// authenticate the user
	// look up the user by email (provided by the JSON payload)
	user, err := app.models.User.GetByEmail(creds.Username)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid username or password"))
		return
	}

	// validate the user's provided password
	validPassword, err := user.PasswordMatches(creds.Password)
	if !validPassword || err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid username or password"))
		return
	}

	// if user is valid, generate a token
	token, err := app.models.Token.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// save token to database
	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}

	// send back response

	payload := jsonResponse{
		Error:   false,
		Message: "logged in",
		Data: envelope{
			"token": token,
			"user":  user,
		},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, err)
		return
	}
}

// Logout logs out a user.
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("invalid json"))
		return
	}

	err = app.models.Token.DeleteByToken(requestPayload.Token)
	if err != nil {
		app.errorLog.Println(err)
		app.errorJSON(w, errors.New("failed to delete token"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged out",
	}

	app.writeJSON(w, http.StatusOK, payload)
}
