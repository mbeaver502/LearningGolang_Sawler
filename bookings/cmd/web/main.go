package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/handlers"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/render"
)

const PORT_NUMBER = ":8080"

var app *config.AppConfig

// main is the program entrypoint.
func main() {
	app = setupAppConfig()

	setupSession(app)
	setupHandlers(app)
	setupTemplates(app)

	// Launch a server and start serving requests on localhost:PORT_NUMBER
	log.Printf("Starting server on port %v\n", PORT_NUMBER)

	srv := &http.Server{
		Addr:    PORT_NUMBER,
		Handler: routes(app),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func setupAppConfig() *config.AppConfig {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = tc
	app.UseCache = false     // change to true when in prod
	app.InProduction = false // change to true when in prod

	return &app
}

func setupSession(a *config.AppConfig) {
	a.Session = scs.New()
	a.Session.Lifetime = 24 * time.Hour
	a.Session.Cookie.Persist = true
	a.Session.Cookie.SameSite = http.SameSiteLaxMode
	a.Session.Cookie.Secure = a.InProduction
}

func setupHandlers(a *config.AppConfig) {
	handlers.NewHandlers(handlers.NewRepo(a))
}

func setupTemplates(a *config.AppConfig) {
	render.NewTemplates(a)
}
