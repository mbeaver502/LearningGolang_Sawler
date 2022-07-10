package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	setupMiddleware(mux)
	setupRoutes(mux)
	setupFileserver(mux)

	return mux
}

func setupMiddleware(m *chi.Mux) {
	m.Use(middleware.Recoverer)
	m.Use(NoSurf)
	m.Use(SessionLoad)
}

func setupRoutes(m *chi.Mux) {
	m.Get("/", handlers.Repo.Home)
	m.Get("/about", handlers.Repo.About)
	m.Get("/generals-quarters", handlers.Repo.Generals)
	m.Get("/majors-suite", handlers.Repo.Majors)
	m.Get("/make-reservation", handlers.Repo.Reservation)
	m.Get("/search-availability", handlers.Repo.Availability)
	m.Post("/search-availability", handlers.Repo.PostAvailability)
	m.Get("/contact", handlers.Repo.Contact)
}

func setupFileserver(m *chi.Mux) {
	// tell chi how to serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	m.Handle("/static/*", http.StripPrefix("/static", fileServer))
}
