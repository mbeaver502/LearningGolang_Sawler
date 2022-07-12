package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/driver"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/handlers"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/helpers"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/render"
)

const PORT_NUMBER = ":8080"

var app *config.AppConfig

// main is the program entrypoint.
func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}

	db := setupDatabase(app)
	defer db.SQL.Close()

	setupLogging(app)
	setupHelpers(app)
	setupSession(app)
	setupHandlers(app, db)
	setupTemplates(app)

	srv := &http.Server{
		Addr:    PORT_NUMBER,
		Handler: routes(app),
	}

	// Launch a server and start serving requests on localhost:PORT_NUMBER
	log.Printf("Starting server on %v\n", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var err error

	app, err = setupAppConfig()
	if err != nil {
		return err
	}

	return nil
}

func setupAppConfig() (*config.AppConfig, error) {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln(err)
		return &app, err
	}

	app.TemplateCache = tc
	app.UseCache = false     // change to true when in prod
	app.InProduction = false // change to true when in prod

	return &app, nil
}

func setupSession(a *config.AppConfig) {
	// what we'll be putting into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	a.Session = scs.New()
	a.Session.Lifetime = 24 * time.Hour
	a.Session.Cookie.Persist = true
	a.Session.Cookie.SameSite = http.SameSiteLaxMode
	a.Session.Cookie.Secure = a.InProduction
}

func setupHandlers(a *config.AppConfig, db *driver.DB) {
	handlers.NewHandlers(handlers.NewRepo(a, db))
}

func setupTemplates(a *config.AppConfig) {
	render.NewRenderer(a)
}

func setupLogging(a *config.AppConfig) {
	a.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	a.ErrorLog = log.New(os.Stdout, "ERR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func setupHelpers(a *config.AppConfig) {
	helpers.NewHelpers(a)
}

func setupDatabase(a *config.AppConfig) *driver.DB {
	log.Println("Connecting to database...")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=password")
	if err != nil {
		log.Fatalln("Failed to connect to database", err)
	}

	log.Println("Connected to database.")

	return db
}
