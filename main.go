package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
)

func main() {
	// Initialize viper to read the configuration file
	// read the einvironment variable stage
	// if not set, set it to empty string
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = ""
	} else {
		stage = "-" + stage
	}

	viper.SetConfigName("config" + stage)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Get configuration values
	port := viper.GetString("server.port")
	requests := viper.GetInt("rate_limit.requests")
	duration := viper.GetDuration("rate_limit.duration")

	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(httprate.Limit(
		requests, // requests
		duration, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Define page routes
	registerPageRoutes(r)

	// Register the API
	registerApi(r)

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
