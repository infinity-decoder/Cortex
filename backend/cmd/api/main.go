package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/infinity-decoder/cortex-backend/internal/auth"
	"github.com/infinity-decoder/cortex-backend/internal/persistence"
	"github.com/infinity-decoder/cortex-backend/internal/scanner"
	"github.com/infinity-decoder/cortex-backend/internal/scheduler"
	"github.com/infinity-decoder/cortex-backend/pkg/db"
)

func main() {
	// Load environment variables from .env if it exists
	godotenv.Load("../.env")

	// Initialize Database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Dependencies
	repo := persistence.NewRepository(database)
	orch := scanner.NewOrchestrator(repo)
	srv := &Server{Repo: repo, Orchestrator: orch}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize & Start Scheduler
	schedule := scheduler.NewScheduler(repo, orch, 24*time.Hour)
	go schedule.Start(ctx)

	r := chi.NewRouter()

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://cortex.security"},
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
		
		// Public Auth Routes
		r.Post("/auth/register", srv.handleRegister)
		r.Post("/auth/login", srv.handleLogin)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware)
			
			r.Post("/scan", srv.handleScan)
			r.Post("/domains/verify", srv.handleVerify)
			r.Get("/stats", srv.handleStats)
			r.Get("/assets", srv.handleGetAssets)
			r.Get("/services", srv.handleGetServices)
			r.Get("/services", srv.handleGetServices)
			r.Get("/findings", srv.handleGetFindings)
			r.Get("/domains", srv.handleGetDomains)

			// Billing Routes
			r.Get("/billing/plan", srv.handleGetPlan)
			r.Put("/billing/plan", srv.handleUpdatePlan)
		})
	})

	port := 8080
	fmt.Printf("Starting Cortex API on :%d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatal(err)
	}
}
