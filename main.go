package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(httprate.Limit(
		10,             // requests
		10*time.Second, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Define page routes
	registerPageRoutes(r)

	// Register the API
	registerApi(r)

	// Get port from environment variable (default to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
