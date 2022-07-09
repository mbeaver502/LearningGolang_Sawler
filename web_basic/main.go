package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a request handler on the main root (e.g., https://localhost:8080/)
	// Note that our handler func is an anonymous function that takes in a writer and *request
	// In this example, we're simply printing out "Hello World!" for every request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Hello World!")

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Number of bytes written: %d; Error: %v\n", n, err)
	})

	// Launch a server and start serving requests on localhost:8080
	http.ListenAndServe(":8080", nil)
}
