package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/10minutemessage/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var messages = cache.Cache{}
var validate = validator.New()
var config Config

type TextMessage struct {
	Text string `json:"text" validate:"required,min=1,max=1000"`
}

func init() {
	// Load configuration
	config = LoadConfig()

	// Update validation rules based on config values
	validate.RegisterValidation("min", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= config.Message.MinLength
	})
	validate.RegisterValidation("max", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) <= config.Message.MaxLength
	})
}

func encode(w http.ResponseWriter, r *http.Request) {
	var textMessage TextMessage

	if r.Body == nil {
		http.Error(w, "No body provided", http.StatusBadRequest)
		return
	}

	if err := render.DecodeJSON(r.Body, &textMessage); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate the message
	if err := validate.Struct(textMessage); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	// Sanitize the message
	textMessage.Text = strings.TrimSpace(textMessage.Text)

	code := uuid.NewString()
	code = strings.Replace(code, "-", "", -1)
	messages.Set(code, textMessage.Text, config.Message.Expiration)

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

	// Sanitize the code
	code = strings.TrimSpace(code)

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
