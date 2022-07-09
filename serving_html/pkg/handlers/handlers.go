package handlers

import (
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/render"
)

// Home is the home page handler.
// Any kind of handler func *must* take in an http.ResponseWriter and a *http.Request!
func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

// About is the about page handler.
func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}

// registerHandlers registers all the handlers for the page routes.
func RegisterHandlers() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
}
