package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/handlers"
)

func routes(ac *config.AppConfig) http.Handler {
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
	m.Get("/contact", handlers.Repo.Contact)

	m.Get("/generals-quarters", handlers.Repo.Generals)
	m.Get("/majors-suite", handlers.Repo.Majors)

	m.Get("/make-reservation", handlers.Repo.Reservation)
	m.Post("/make-reservation", handlers.Repo.PostReservation)
	m.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	m.Get("/search-availability", handlers.Repo.Availability)
	m.Post("/search-availability", handlers.Repo.PostAvailability)
	m.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	m.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	m.Get("/book-room", handlers.Repo.BookRoom)

	m.Get("/user/login", handlers.Repo.ShowLogin)
	m.Post("/user/login", handlers.Repo.PostShowLogin)
	m.Get("/user/logout", handlers.Repo.Logout)

	// Allow only authenticated users to access /admin/<route>
	m.Route("/admin", func(m chi.Router) {
		//m.Use(Auth)

		m.Get("/dashboard", handlers.Repo.AdminDashboard)

		m.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		m.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		m.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)

		m.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		m.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
		m.Get("/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservation)
	})
}

func setupFileserver(m *chi.Mux) {
	// tell chi how to serve static files
	fileServer := http.FileServer(http.Dir("./static/"))
	m.Handle("/static/*", http.StripPrefix("/static", fileServer))
}
