package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApplication_AllUsers(t *testing.T) {
	// create some mock rows, and add one row
	var mockedRows = mockedDB.NewRows([]string{"id", "email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at", "has_token"})
	mockedRows.AddRow("1", "me@here.com", "John", "Doe", "abc123", "1", time.Now(), time.Now(), "1")

	// tell mocked DB what queries we expect
	mockedDB.ExpectQuery("select \\\\* ").WillReturnRows(mockedRows)

	// create a test recorder
	rr := httptest.NewRecorder()

	// create a new request
	req, _ := http.NewRequest("POST", "/admin/users", nil)

	// call the handler
	handler := http.HandlerFunc(testApp.AllUsers)
	handler.ServeHTTP(rr, req)

	// check for expected status code
	if rr.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rr.Code, http.StatusOK)
	}
}
