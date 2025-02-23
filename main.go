package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// create map to store the messages
var messages = make(map[string]string)

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
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/", filesDir)

	// regsiter the API
	registerApi(r)

	// Start the server
	http.ListenAndServe(":8080", r)
}

func registerApi(r *chi.Mux) {
	r.Post("/encode", encode)
	r.Route("/decode", func(r chi.Router) {
		r.Get("/{code}", decode)
	})
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
	messages[code] = textMessage.Text

	var responseMessage string
	responseMessage = "decode/" + code

	render.JSON(w, r, responseMessage)
}

func decode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	text, ok := messages[code]
	if !ok {
		http.Error(w, "Code not found", http.StatusNotFound)
		return
	}
	render.JSON(w, r, text)
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
