package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/10minutemessage/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var messages = cache.Cache{}

type TextMessage struct {
	Text string `json:"text"`
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
