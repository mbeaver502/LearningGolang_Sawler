package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_Routes_Exist(t *testing.T) {
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	// these routes must exist (incomplete list...)
	expectedRoutes := []string{
		"/users/login",
		"/users/logout",
		"/admin/users/get/{id}",
		"/admin/users/save",
		"/admin/users",
		"/admin/users/delete",
	}

	for _, route := range expectedRoutes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	// assume that the route does not exist
	found := false

	// walk through all registered routes
	// look for the desired route
	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// if we find route we're looking for,
		// set found to true
		if route == foundRoute {
			found = true
		}
		return nil
	})

	// fire an error if route not found
	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
