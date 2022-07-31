package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	// what we expect to receive in our request
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	// attempt to parse the request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// create our message
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// try to send the email
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// write back to the caller
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("email sent to %s", requestPayload.To),
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
