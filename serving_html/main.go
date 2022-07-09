package main

import (
	"log"
	"net/http"
)

const PORT_NUMBER = ":8080"

// main is the program entrypoint.
func main() {
	registerHandlers()

	log.Printf("Starting server on port %v\n", PORT_NUMBER)

	// Launch a server and start serving requests on localhost:8080
	http.ListenAndServe(PORT_NUMBER, nil)
}
