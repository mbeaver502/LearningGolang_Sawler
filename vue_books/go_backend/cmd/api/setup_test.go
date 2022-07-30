package main

import (
	"books_backend/internal/data"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var testApp application
var mockedDB sqlmock.Sqlmock

func TestMain(m *testing.M) {
	testDB, myMock, _ := sqlmock.New()
	mockedDB = myMock
	defer testDB.Close()

	testApp = application{
		config:      config{},
		models:      data.New(testDB),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:    log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		environment: "development",
	}

	os.Exit(m.Run())
}
