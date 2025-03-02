package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
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
