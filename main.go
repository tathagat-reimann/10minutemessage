package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/10minutemessage/cache"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var messages = cache.Cache{}

type TextMessage struct {
	Text string `json:"text"`
}

func main() {
	fmt.Println("Starting 10minutemessage!")

	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// Define page routes
	registerPageRoutes(r)

	// regsiter the API
	registerApi(r)

	// Start the server
	http.ListenAndServe(":8080", r)
}

func registerPageRoutes(r *chi.Mux) {
	r.Get("/", handleMainUrl)
	r.Get("/e", handleEncode)
	r.Get("/d", handleDecodeText)
}

func registerApi(r *chi.Mux) {

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/encode", encode)
		r.Route("/decode", func(r chi.Router) {
			r.Get("/{code}", decode)
		})
	})
}

func handleMainUrl(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/e", http.StatusFound)
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/encode/index.html")
}

func handleDecodeText(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/decode/index.html")
}

func encode(w http.ResponseWriter, r *http.Request) {
	var textMessage TextMessage
	if err := render.DecodeJSON(r.Body, &textMessage); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if textMessage.Text == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}
	code := uuid.New().String()
	// susbstitute the "-" with ""
	code = strings.Replace(code, "-", "", -1)
	messages.Set(code, textMessage.Text, 10*time.Minute)

	var responseMessage string
	responseMessage = "/d?code=" + code

	render.JSON(w, r, responseMessage)
}

func decode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	text, ok := messages.Get(code)

	if !ok {
		http.Error(w, "Code not found", http.StatusNotFound)
		return
	}
	render.PlainText(w, r, text)
}
