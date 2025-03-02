package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/10minutemessage/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var messages = cache.Cache{}

type TextMessage struct {
	Text string `json:"text"`
}

func main() {
	log.Println("Starting 10MinuteMessage!")

	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

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

func registerPageRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/e", http.StatusFound)
	})
	r.Get("/e", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/encode/index.html")
	})
	r.Get("/d", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/decode/index.html")
	})
}

func registerApi(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/encode", encode)
		r.Route("/decode", func(r chi.Router) {
			r.Get("/{code}", decode)
		})
	})
}

func encode(w http.ResponseWriter, r *http.Request) {
	var textMessage TextMessage
	if err := render.DecodeJSON(r.Body, &textMessage); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if textMessage.Text == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	code := uuid.NewString()
	// susbstitute the "-" with ""
	code = strings.Replace(code, "-", "", -1)
	messages.Set(code, textMessage.Text, 10*time.Minute)

	response := map[string]string{
		"message": "Message encoded successfully",
		"url":     "/d?code=" + code,
	}

	render.JSON(w, r, response)
}

func decode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"text": "",
	}

	text, ok := messages.Get(code)
	if !ok {
		http.Error(w, "Code not found", http.StatusNotFound)
		return
	}

	response["text"] = text
	render.JSON(w, r, response)
}
