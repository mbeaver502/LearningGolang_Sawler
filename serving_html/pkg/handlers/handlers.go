package handlers

import (
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/models"
	"github.com/mbeaver502/LearningGolang_Sawler/serving_html/pkg/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler.
// Any kind of handler func *must* take in an http.ResponseWriter and a *http.Request!
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	render.RenderTemplate_v3(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
