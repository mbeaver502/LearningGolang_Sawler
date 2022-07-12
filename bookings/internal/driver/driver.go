package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	MAX_OPEN_DB_CONN = 10
	MAX_IDLE_DB_CONN = 5
	MAX_DB_LIFETIME  = 5 * time.Minute
	POSTGRES_DRIVER  = "pgx"
)

// ConnectSQL creates a database pool for Postgres.
func ConnectSQL(connectionString string) (*DB, error) {
	d, err := NewDatabase(connectionString)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(MAX_OPEN_DB_CONN)
	d.SetMaxIdleConns(MAX_IDLE_DB_CONN)
	d.SetConnMaxLifetime(MAX_DB_LIFETIME)

	dbConn.SQL = d

	// test the DB connection one more time
	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping the given database.
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	return nil
}

// NewDatabase creates a new database connection for the given connection string.
func NewDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open(POSTGRES_DRIVER, connectionString)
	if err != nil {
		return nil, err
	}

	// test the database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
