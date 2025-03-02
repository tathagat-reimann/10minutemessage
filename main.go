package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	// Load configuration
	config := LoadConfig()

	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(httprate.Limit(
		config.Requests, // requests
		config.Duration, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Define page routes
	registerPageRoutes(r)

	// Register the API
	registerApi(r)

	log.Printf("Server running on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}
