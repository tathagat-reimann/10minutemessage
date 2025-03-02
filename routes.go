package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

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
