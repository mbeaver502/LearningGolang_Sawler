package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT_NUMBER = ":8080"

// Home is the home page handler.
// Any kind of handler func *must* take in an http.ResponseWriter and a *http.Request!
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page.")
}

// About is the about page handler.
func About(w http.ResponseWriter, r *http.Request) {
	x, y := 2, 2
	sum := addValues(x, y)

	fmt.Fprintf(w, "This is the about page. And %v + %v is %v", x, y, sum)
}

// Hello is the hello page handler.
func Hello(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "Hello World!")

	if err != nil {
		log.Println("hello error:", err)
	}

	log.Printf("Number of bytes written: %d; Error: %v\n", n, err)
}

// addValues adds two integers and returns the sum.
// Notice that this is a normal function, not a request handler!
// Our addValues func will be unexported, so only package main can access it
func addValues(x, y int) int {
	return x + y
}

// main is the program entrypoint.
func main() {
	registerHandlers()

	log.Printf("Starting server on port %v\n", PORT_NUMBER)

	// Launch a server and start serving requests on localhost:8080
	http.ListenAndServe(PORT_NUMBER, nil)
}

// registerHandlers registers all the handlers for the page routes.
func registerHandlers() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/hello", Hello)
}
