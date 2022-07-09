package main

import (
	"log"
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/handlers"
)

const PORT_NUMBER = ":8080"

// main is the program entrypoint.
func main() {
	handlers.RegisterHandlers()

	log.Printf("Starting server on port %v\n", PORT_NUMBER)

	// Launch a server and start serving requests on localhost:8080
	http.ListenAndServe(PORT_NUMBER, nil)
}
