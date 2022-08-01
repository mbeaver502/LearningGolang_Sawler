package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}

	for _, r := range routes {
		routeExists(t, chiRoutes, r)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("want %s; did not find", route)
	}
}
