package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPageRoutes(t *testing.T) {
	r := chi.NewRouter()
	registerPageRoutes(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "/", http.StatusFound},
		{"GET", "/e", http.StatusOK},
		{"GET", "/d", http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestRegisterApiNilBody(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/encode", http.StatusBadRequest},       // No body provided, should return 400
		{"GET", "/api/decode/somecode", http.StatusNotFound}, // Assuming "somecode" does not exist
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestRegisterApiBadBody(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/encode", http.StatusBadRequest},       // No body provided, should return 400
		{"GET", "/api/decode/somecode", http.StatusNotFound}, // Assuming "somecode" does not exist
	}

	userJSON := `{"This is not json"}`
	body := strings.NewReader(userJSON)

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, body)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestEncodeEmptyText(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/encode", http.StatusBadRequest},
	}

	userJSON := `{"text":""}`
	body := strings.NewReader(userJSON)

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, body)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestEncode(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/api/encode", http.StatusOK},
	}

	userJSON := `{"text":"asdf"}`
	body := strings.NewReader(userJSON)

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, body)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestDecodeNotFound(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "/api/decode", http.StatusNotFound},
	}

	userJSON := `{"text":"asdf"}`
	body := strings.NewReader(userJSON)

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, body)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}

func TestDecodeFound(t *testing.T) {
	r := chi.NewRouter()
	registerApi(r)

	// preparing the code
	userJSON := `{"text":"asdf"}`
	body := strings.NewReader(userJSON)
	req, _ := http.NewRequest("POST", "/api/encode", body)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	// fmt.Println(rr.Body.String())
	responseMap := make(map[string]string)
	render.DecodeJSON(rr.Body, &responseMap)
	url := responseMap["url"]
	// get code which is the string after the first "="
	code := url[strings.Index(url, "=")+1:]

	// fmt.Println("Code: ", code)

	req, err := http.NewRequest("GET", "/api/decode/"+code, nil)
	assert.NoError(t, err)

	// rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	responseMap2 := make(map[string]string)
	render.DecodeJSON(rr.Body, &responseMap2)
	text := responseMap2["text"]
	assert.Equal(t, "asdf", text)
}
