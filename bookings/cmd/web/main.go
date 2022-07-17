package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
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
	// read flags
	inProduction := flag.Bool("production", true, "Application is in Production")
	useCache := flag.Bool("cache", true, "Use cached templates")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		log.Fatalln(errors.New("missing required flags"))
	}

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
	defer close(app.MailChan)
	listenForMail()

	app.InProduction = *inProduction
	app.UseCache = *useCache

	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db := setupDatabase(app, connString)
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

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	return &app, nil
}

func setupSession(a *config.AppConfig) {
	// what we'll be putting into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

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

func setupDatabase(a *config.AppConfig, connString string) *driver.DB {
	log.Println("Connecting to database...")

	db, err := driver.ConnectSQL(connString)
	if err != nil {
		log.Fatalln("Failed to connect to database", err)
	}

	log.Println("Connected to database.")

	return db
}
