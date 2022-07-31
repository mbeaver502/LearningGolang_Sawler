package main

import (
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RequestPayload is the standard, expected JSON structure.
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
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

// MailPayload represents the JSOn for a mail request.
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
		//app.logItem(w, requestPayload.Log)
		// app.logEventViaRabbit(w, requestPayload.Log)
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
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

/* func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
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
} */

func (app *Config) sendMail(w http.ResponseWriter, mail MailPayload) {
	jsonData, _ := json.Marshal(mail)

	// call the Mailer service
	mailServiceURL := "http://mailer-service/send" // mailer-service is the name inside Docker
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// ensure we got the right status code
	if resp.StatusCode != http.StatusAccepted {
		log.Printf("invalid status code: %d", resp.StatusCode)
		app.errorJSON(w, fmt.Errorf("invalid status code: %d", resp.StatusCode))
		return
	}

	// write a response back to the front-end
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("mail sent to %s", mail.To),
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

/*
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged via rabbitmq",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name string, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return err
	}

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
*/

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	// logger-service is the name of the Logger service in Docker
	// the Logger service program is set up to listen for RPC on Port 5001
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// what we'll send to the RPC listener
	// rpcPayload := RPCPayload{
	// 	Name: l.Name,
	// 	Data: l.Data,
	// }
	rpcPayload := RPCPayload(l)

	var result string

	// the serviceMethod name must match exactly
	// the serviceMethod must be exported
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// LogViaGRPC will log something to Mongo via a gRPC request.
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// logger-service is the name of the service running inside Docker
	// Logger is listening for gRPC on port 50001
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// set up a new client that is designed specifically for logs
	// based on the logs.proto
	c := logs.NewLogServiceClient(conn)

	// get a context for our gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// try to write to the log by executing a gRPC
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	// write back a response
	payload := jsonResponse{
		Error:   false,
		Message: "logged via gRPC",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
