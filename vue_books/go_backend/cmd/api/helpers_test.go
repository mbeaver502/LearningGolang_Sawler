package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication_readJSON(t *testing.T) {
	// create sample JSON
	sampleJSON := map[string]interface{}{
		"foo": "bar",
	}

	body, _ := json.Marshal(sampleJSON)

	// declare var to hold read values
	var decodedJSON struct {
		Foo string `json:"foo"`
	}

	// create a request
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Log(err)
	}

	// create a response recorder
	rr := httptest.NewRecorder()
	defer req.Body.Close()

	// try to read JSON
	err = testApp.readJSON(rr, req, &decodedJSON)
	if err != nil {
		t.Error(err)
	}

	if decodedJSON.Foo != sampleJSON["foo"] {
		t.Errorf("got %s; want %s", decodedJSON.Foo, sampleJSON["foo"])
	}
}

func TestApplication_writeJSON(t *testing.T) {
	// create a response recorder
	rr := httptest.NewRecorder()

	payload := jsonResponse{
		Error:   false,
		Message: "foo",
		Data:    nil,
	}

	// create some test headers
	headers := make(http.Header)
	headers.Add("Foo", "Bar")

	// try to write JSON
	err := testApp.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Error(err)
	}
}

func TestApplication_errorJSON(t *testing.T) {
	// create a response recorder
	rr := httptest.NewRecorder()

	// test writing simple error JSON
	err := testApp.errorJSON(rr, errors.New("test error"))
	if err != nil {
		t.Error(err)
	}

	testJSONPayload(t, rr)

	// specifically caught Postgres error codes
	errSlice := []string{
		"(SQLSTATE 23505)",
		"(SQLSTATE 22001)",
		"(SQLSTATE 23503)",
	}

	for _, x := range errSlice {
		customErr := testApp.errorJSON(rr, errors.New(x), http.StatusUnauthorized)
		if customErr != nil {
			t.Error(err)
		}

		testJSONPayload(t, rr)
	}
}

func testJSONPayload(t *testing.T, rr *httptest.ResponseRecorder) {
	var requestPayload jsonResponse

	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&requestPayload)
	if err != nil {
		t.Error(err)
	}

	// we expect errorJSON's payload to have Error = true
	if !requestPayload.Error {
		t.Errorf("got error = false in response; want error = true")
	}
}
