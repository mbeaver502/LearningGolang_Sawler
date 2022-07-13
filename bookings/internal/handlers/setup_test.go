package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

const (
	TEMPLATES_DIRECTORY = "./templates"
	TEMPLATE_FILE       = "*.page.tmpl"
	LAYOUT_FILE         = "*.layout.tmpl"
)

var pathToTemplates string = "./../../templates"

func TestMain(m *testing.M) {
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = tc
	app.UseCache = true      // change to true when in prod
	app.InProduction = false // change to true when in prod

	// what we'll be putting into the session
	gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	setupLogging(&app)
	NewHandlers(NewTestRepo(&app))
	render.NewRenderer(&app)
	routes(&app)

	os.Exit(m.Run())
}

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	setupMiddleware(mux)
	setupRoutes(mux)
	setupFileserver(mux)

	return mux
}

func setupMiddleware(m *chi.Mux) {
	m.Use(middleware.Recoverer)
	//m.Use(NoSurf)
	m.Use(SessionLoad)
}

func setupRoutes(m *chi.Mux) {
	m.Get("/", Repo.Home)
	m.Get("/about", Repo.About)
	m.Get("/contact", Repo.Contact)

	m.Get("/generals-quarters", Repo.Generals)
	m.Get("/majors-suite", Repo.Majors)

	m.Get("/make-reservation", Repo.Reservation)
	m.Post("/make-reservation", Repo.PostReservation)
	m.Get("/reservation-summary", Repo.ReservationSummary)

	m.Get("/search-availability", Repo.Availability)
	m.Post("/search-availability", Repo.PostAvailability)
	m.Post("/search-availability-json", Repo.AvailabilityJSON)
}

func setupFileserver(m *chi.Mux) {
	// tell chi how to serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	m.Handle("/static/*", http.StripPrefix("/static", fileServer))
}

// NoSurf adds CSRF protection to all POST requests.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request.
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

// createTemplateCache creates a template cache
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	templateFiles := fmt.Sprintf("%s/%s", pathToTemplates, TEMPLATE_FILE) //strings.Join([]string{TEMPLATES_DIRECTORY, TEMPLATE_FILE}, "")
	layoutFiles := fmt.Sprintf("%s/%s", pathToTemplates, LAYOUT_FILE)     // strings.Join([]string{TEMPLATES_DIRECTORY, LAYOUT_FILE}, "")

	// get all our page template files
	pages, err := filepath.Glob(templateFiles)
	if err != nil {
		log.Println(err)
		return cache, err
	}

	// range over all page template files
	for _, page := range pages {
		// get the name of the file itself and parse it as a template with that name
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// get all our layout template files
		layouts, err := filepath.Glob(layoutFiles)
		if err != nil {
			log.Println(err)
			return cache, err
		}

		// associate any layouts with templates that require them
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(layoutFiles)
			if err != nil {
				log.Println(err)
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

func setupLogging(a *config.AppConfig) {
	a.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	a.ErrorLog = log.New(os.Stdout, "ERR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
