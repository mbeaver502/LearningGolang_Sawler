package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
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

// renderTemplate renders an HTML template to the given http.ResponseWriter.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	path := strings.Join([]string{"./templates/", tmpl}, "")

	parsedTemplate, _ := template.ParseFiles(path)
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		log.Println("renderTemplate error:", err)
		return
	}
}
