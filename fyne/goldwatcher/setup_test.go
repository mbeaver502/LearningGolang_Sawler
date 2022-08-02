package main

import (
	"net/http"
	"os"
	"testing"
)

var testApp Config

func TestMain(m *testing.M) {

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
