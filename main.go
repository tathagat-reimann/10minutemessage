package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	// Deliver index.html
	workDir, _ := os.Getwd()
	encodeFilesDir := http.Dir(filepath.Join(workDir, "static/encode"))
	FileServer(r, "/encode", encodeFilesDir)
	decodeFilesDir := http.Dir(filepath.Join(workDir, "static/encode"))
	FileServer(r, "/decode", decodeFilesDir)

	// regsiter the API
	registerApi(r)

	// Start the server
	http.ListenAndServe(":8080", r)
}

func registerApi(r *chi.Mux) {
	// redirect base URL to /encode
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/encode", http.StatusFound)
	})

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/encode", encode)
		r.Route("/decode", func(r chi.Router) {
			r.Get("/{code}", decode)
		})
	})

	/*
		r.Route("/blogs", func(r chi.Router) {
			r.Get("/", getAllBlogs)
			r.Get("/{id}", getBlog)
			r.Post("/", createBlog)
		})
	*/
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
	//[code] = textMessage.Text

	var responseMessage string
	responseMessage = "/decode/" + code

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

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
