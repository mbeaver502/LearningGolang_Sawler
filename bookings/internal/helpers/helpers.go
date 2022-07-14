package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers.
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated determines if a user is authenticated by checking session variable exists.
func IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}
