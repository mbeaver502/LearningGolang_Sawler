package main

import (
	"net/http"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/mbeaver502/LearningGolang_Sawler/bookings/internal/config"
)

func TestNosurf(t *testing.T) {
	var mh mockHandler
	app = &config.AppConfig{
		InProduction: false,
	}

	h := NoSurf(&mh)

	switch v := h.(type) {
	case http.Handler:
		// expected
	default:
		t.Errorf("type is not http.Handler, but is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var mh mockHandler
	app = &config.AppConfig{
		Session: scs.New(),
	}

	h := SessionLoad(&mh)

	switch v := h.(type) {
	case http.Handler:
		// expected
	default:
		t.Errorf("type is not http.Handler, but is %T", v)
	}
}
