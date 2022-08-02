package main

import (
	"bytes"
	"goldwatcher/repository"
	"io"
	"net/http"
	"os"
	"testing"

	"fyne.io/fyne/v2/test"
)

var testApp Config

var client = NewTestClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer([]byte(jsonToReturn))),
	}
})

func TestMain(m *testing.M) {
	testApp.App = test.NewApp()
	testApp.HTTPClient = client
	testApp.MainWindow = testApp.App.NewWindow("")
	testApp.DB = repository.NewTestRepository()

	os.Exit(m.Run())
}

var jsonToReturn = `{
	"ts": 1659397837572,
	"tsj": 1659397835878,
	"date": "Aug 1st 2022, 07:50:35 pm NY",
	"items": [
	  {
		"curr": "USD",
		"xauPrice": 1772.005,
		"xagPrice": 20.3385,
		"chgXau": 9.2726,
		"chgXag": 0.1161,
		"pcXau": 0.526,
		"pcXag": 0.5741,
		"xauClose": 1762.73243,
		"xagClose": 20.22242
	  }
	]
  }`

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
