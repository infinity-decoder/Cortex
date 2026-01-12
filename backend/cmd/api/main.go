package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Cortex API is healthy"))
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Cortex API v1"))
		})
		
		r.Post("/scan", handleScan)
	})

	port := 8080
	fmt.Printf("Starting Cortex API on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
