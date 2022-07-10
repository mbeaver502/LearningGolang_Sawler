package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/models"
)

type mockWriter struct{}

func (w *mockWriter) Header() http.Header         { return http.Header{} }
func (w *mockWriter) Write(b []byte) (int, error) { return len(b), nil }
func (w *mockWriter) WriteHeader(statusCode int)  {}

func TestMain(m *testing.M) {
	testApp, _ := setupAppConfig()
	setupSession(testApp)
	setupLogging(testApp)

	// set the app in main render package
	app = testApp

	os.Exit(m.Run())
}

func setupAppConfig() (*config.AppConfig, error) {
	var app config.AppConfig

	app.UseCache = false     // change to true when in prod
	app.InProduction = false // change to true when in prod

	return &app, nil
}

func setupSession(a *config.AppConfig) {
	// what we'll be putting into the session
	gob.Register(models.Reservation{})

	a.Session = scs.New()
	a.Session.Lifetime = 24 * time.Hour
	a.Session.Cookie.Persist = true
	a.Session.Cookie.SameSite = http.SameSiteLaxMode
	a.Session.Cookie.Secure = a.InProduction
}

func setupLogging(a *config.AppConfig) {
	a.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	a.ErrorLog = log.New(os.Stdout, "ERR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
