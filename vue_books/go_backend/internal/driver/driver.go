package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	maxOpenDBConn = 5
	maxIdleDBConn = 5
	maxDBLifetime = 5 * time.Minute
)

func ConnectPostgres(dsn string) (*DB, error) {
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	err = testDB(err, d)
	dbConn.SQL = d

	return dbConn, err
}

func testDB(err error, d *sql.DB) error {
	err = d.Ping()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("*** Pinged database successfully! ***")

	return err
}
