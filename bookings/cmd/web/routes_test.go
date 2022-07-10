package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	app, _ := setupAppConfig()

	setupSession(app)
	setupHandlers(app)
	setupTemplates(app)

	mux := routes(app)

	switch v := mux.(type) {
	case *chi.Mux:
		// expected
	default:
		t.Errorf("expected *chi.Mux, got %T", v)
	}
}
