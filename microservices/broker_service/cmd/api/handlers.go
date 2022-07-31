package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// RequestPayload is the standard, expected JSON structure.
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

// AuthPayload represents the JSON for an authentication request.
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogPayload represents the JSON for a logging request.
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Broker is a sample handler.
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "broker says hello world!",
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the default handler for all requests.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	// attempt to parse the request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// act according to what the user requested
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unrecognized action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some JSON to be sent to the auth service
	jsonData, _ := json.Marshal(a)

	// call the auth service -- using the host service we named in docker-compose
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// make sure we get back the correct status code
	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable to read resp.Body into
	var jsonFromService jsonResponse

	// decode the JSON from the auth service
	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// if auth service said there's an error, tell user
	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "authenticated",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	// turn the log request into JSON we can send to Logger
	jsonData, _ := json.Marshal(entry)

	// call the service
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// ensure we got the right response status code
	if response.StatusCode != http.StatusAccepted {
		log.Printf("invalid status code: %d", response.StatusCode)
		app.errorJSON(w, fmt.Errorf("invalid status code: %d", response.StatusCode))
		return
	}

	// write a response back to the front-end
	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
