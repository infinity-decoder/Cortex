package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"cortex-backend/internal/auth"
	"cortex-backend/internal/ratelimit"
	httpsmiddleware "cortex-backend/internal/middleware"
	"cortex-backend/internal/persistence"
	"cortex-backend/internal/queue"
	"cortex-backend/internal/scanner"
	"cortex-backend/internal/scheduler"
	"cortex-backend/pkg/db"
)

func main() {
	// Load environment variables from .env if it exists
	godotenv.Load("../.env")

	// Validate JWT Secret Key
	if err := auth.ValidateJWTSecret(); err != nil {
		log.Fatalf("JWT Secret Key validation failed: %v", err)
	}

	// Initialize Database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Dependencies
	repo := persistence.NewRepository(database)
	orch := scanner.NewOrchestrator(repo)
	jobQueue := queue.GetQueue()
	
	// Start scan worker
	worker := queue.NewWorker(jobQueue, orch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.Start(ctx)
	
	srv := &Server{Repo: repo, Orchestrator: orch, Queue: jobQueue}

	// Initialize & Start Scheduler
	schedule := scheduler.NewScheduler(repo, orch, 24*time.Hour)
	go schedule.Start(ctx)

	r := chi.NewRouter()

	// CORS Configuration from environment
	corsOriginsStr := os.Getenv("CORS_ORIGINS")
	var allowedOrigins []string
	if corsOriginsStr != "" {
		allowedOrigins = strings.Split(corsOriginsStr, ",")
		// Trim whitespace
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	} else {
		// Default origins for development
		allowedOrigins = []string{"http://localhost:3000"}
	}

	// CORS Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
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
	
	// HTTPS Enforcement (only in production)
	r.Use(httpsmiddleware.HTTPSRedirect)
	r.Use(httpsmiddleware.HSTS)
	
	// Global rate limiting (loose limits)
	globalLimiter := ratelimit.NewPerIPRateLimiter(100.0, 200) // 100 req/sec, burst of 200
	r.Use(globalLimiter.Limit)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Cortex API is healthy"))
	})

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Cortex API v1"))
		})
		
		// Public Auth Routes with stricter rate limiting
		authLimiter := ratelimit.NewPerIPRateLimiter(5.0, 10) // 5 req/sec, burst of 10
		r.With(authLimiter.Limit).Post("/auth/register", srv.handleRegister)
		r.With(authLimiter.Limit).Post("/auth/login", srv.handleLogin)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware)
			
			// Scan endpoint with moderate rate limiting
			scanLimiter := ratelimit.NewPerIPRateLimiter(10.0, 20) // 10 req/sec, burst of 20
			r.With(scanLimiter.Limit).Post("/scan", srv.handleScan)
			r.Post("/domains", srv.handleCreateDomain)
			r.Post("/domains/verify", srv.handleVerify)
			r.Get("/stats", srv.handleStats)
			r.Get("/assets", srv.handleGetAssets)
			r.Get("/services", srv.handleGetServices)
			r.Get("/findings", srv.handleGetFindings)
			r.Get("/domains", srv.handleGetDomains)
			r.Get("/domains/all", srv.handleGetAllDomains)
			r.Get("/scans/status", srv.handleGetScanStatus)

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
