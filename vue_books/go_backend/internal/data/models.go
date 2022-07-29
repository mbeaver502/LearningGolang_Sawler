package data

import (
	"database/sql"
	"time"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

// Models represents available data models.
type Models struct {
	User  User
	Token Token
}

// New sets the database pool and returns available data models.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}
}
