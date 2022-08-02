package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

var testRepo *SQLiteRepository

func TestMain(m *testing.M) {
	path := "./testdata/sql.db"
	_ = os.Remove(path)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	testRepo = NewSQLiteRepository(db)

	os.Exit(m.Run())
}
