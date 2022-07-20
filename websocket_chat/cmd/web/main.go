package main

import (
	"log"
	"net/http"
)

const SERVER_ADDRESS = ":8080"

func main() {
	mux := routes()

	log.Println("Starting server on", SERVER_ADDRESS)

	_ = http.ListenAndServe(SERVER_ADDRESS, mux)
}
