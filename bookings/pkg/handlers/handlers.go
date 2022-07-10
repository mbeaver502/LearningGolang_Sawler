package handlers

import (
	"net/http"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/models"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/render"
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
	// put the user's IP address into our session
	m.App.Session.Put(r.Context(), "remote_ip", r.RemoteAddr)

	render.RenderTemplate_v3(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// get the user's IP address out of the session ("" if it doesn't exist)
	// and put it into our stringmap
	stringMap["remote_ip"] = m.App.Session.GetString(r.Context(), "remote_ip")

	render.RenderTemplate_v3(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form.
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the General's Quarters page.
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Major's Suite page.
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders a page with room availabilities.
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact renders a with contact information.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate_v3(w, "contact.page.tmpl", &models.TemplateData{})
}
