package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Authentication listening on port %s\n", webPort)

	// set up connection to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("failed to connect to postgres!")
	}
	defer conn.Close()

	// set up app config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("*** Successfully pinged database! ***")

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)

		if err != nil {
			log.Println("postgres not yet ready...", err)
			counts++
		} else {
			log.Println("*** Connected to Postgres! ***")
			return conn
		}

		if counts > 10 {
			log.Println("failed connecting for 10 tries...")
			return nil
		}

		log.Println("backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
