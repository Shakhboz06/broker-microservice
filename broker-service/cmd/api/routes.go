package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
)

func(app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//speicify who is allowed
	mux.Use(middleware.Logger)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		ExposedHeaders: []string{"Link"},
		MaxAge: 300,
	}))

	mux.Post("/",app.Handler)

	mux.Post("/log-grpc", app.logWithgRPC)

	mux.Post("/handle", app.HandleSubmissions)

	return mux
}



 