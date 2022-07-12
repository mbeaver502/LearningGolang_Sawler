package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to a database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=postgres password=password")
	if err != nil {
		log.Fatalln("Unable to connect to database", err)
	}
	defer conn.Close()

	log.Println("Connect to database")

	// test the connection
	err = conn.Ping()
	if err != nil {
		log.Fatalln("Cannot ping database", err)
	}

	log.Println("Pinged the database")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatalln("Failed to get rows", err)
	}

	// insert a row
	// $1 and $2 are placeholders for the prepared statement
	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Steve", "Rogers")
	if err != nil {
		log.Fatalln("Error while inserting", err)
	}

	log.Println("Inserted a row")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalln("Failed to get rows", err)
	}

	// update a row
	stmt := `update users set last_name = $1 where id = $2`
	_, err = conn.Exec(stmt, "Brown", 4)
	if err != nil {
		log.Fatalln("Failed to update rows", err)
	}

	log.Println("Updated row(s)")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalln("Failed to get rows", err)
	}

	// get one row by id
	query = `select id, first_name, last_name from users where id = $1`
	row := conn.QueryRow(query, 5)

	var firstName, lastName string
	var id int

	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatalln("Failed to scan one row", err)
	}

	log.Println("QueryRow returns", id, firstName, lastName)

	// delete a row
	stmt = `delete from users where id >= $1`
	_, err = conn.Exec(stmt, 3)
	if err != nil {
		log.Fatalln("Failed to delete", err)
	}

	log.Println("Deleted row(s)")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalln("Failed to get rows", err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println("Query failed", err)
		return err
	}

	// if a query can return more than one row, we must close the query connection
	defer rows.Close()

	var firstName, lastName string
	var id int

	log.Println("------------")

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println("Scan failed", err)
			return err
		}

		log.Println("record is:", id, firstName, lastName)
	}

	// it's good practice to double-check after scanning that there were no errors
	if err = rows.Err(); err != nil {
		log.Fatalln("Error scanning rows", err)
		return err
	}

	log.Println("------------")

	return nil
}
