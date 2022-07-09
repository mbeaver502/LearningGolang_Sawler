package main

import (
	"net/http"
)

// Home is the home page handler.
// Any kind of handler func *must* take in an http.ResponseWriter and a *http.Request!
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.tmpl")
}

// About is the about page handler.
func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.tmpl")
}

// registerHandlers registers all the handlers for the page routes.
func registerHandlers() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
}
