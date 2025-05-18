package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{
		"/authentication",
	}

	for _, r := range routes{
		 routeExists(t, chiRoutes, r)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string){
	found := false
	chi.Walk(routes, func (method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute{
			found = true
		}
		return nil
	})


	if !found{
		t.Errorf("did not find %s in registered routes", route)
	}
}