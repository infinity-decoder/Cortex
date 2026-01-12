package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/infinity-decoder/cortex-backend/pkg/db"
)

func main() {
	// Initialize Database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	r := chi.NewRouter()

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatal(err)
	}
}
