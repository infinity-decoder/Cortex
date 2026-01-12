package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infinity-decoder/cortex-backend/internal/discovery"
	"github.com/infinity-decoder/cortex-backend/internal/persistence"
	"github.com/infinity-decoder/cortex-backend/internal/risk"
	"github.com/infinity-decoder/cortex-backend/internal/scanner"
)

type Server struct {
	Repo         *persistence.Repository
	Orchestrator *scanner.Orchestrator
}

type ScanRequest struct {
	Domain string `json:"domain"`
}

type ScanResponse struct {
	Domain      string             `json:"domain"`
	Assets      []discovery.Result `json:"assets"`
	NewFindings []risk.Exposure    `json:"new_findings"`
	AllFindings []risk.Exposure    `json:"all_findings"`
	Verified    bool               `json:"verified"`
}

type VerifyRequest struct {
	Domain string `json:"domain"`
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Check if Domain is verified
	domain, err := s.Repo.GetDomainByName(ctx, req.Domain)
	if err != nil || !domain.Verified {
		// Auto-initialize if not exists, but block scan
		if err != nil {
			s.Repo.GetOrCreateDomain(ctx, "00000000-0000-0000-0000-000000000000", req.Domain)
		}
		
		http.Error(w, "Domain not verified. Please verify ownership via DNS TXT record first.", http.StatusForbidden)
		return
	}

	// 2. Run Scan via Orchestrator
	result, err := s.Orchestrator.RunScan(ctx, req.Domain, domain.ID.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := ScanResponse{
		Domain:      req.Domain,
		Assets:      result.Assets,
		NewFindings: result.NewFindings,
		AllFindings: result.AllFindings,
		Verified:    true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domain, err := s.Repo.GetDomainByName(ctx, req.Domain)
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	success, err := discovery.VerifyDomain(req.Domain, domain.VerificationToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification check failed: %v", err), http.StatusInternalServerError)
		return
	}

	if success {
		s.Repo.UpdateDomainVerification(ctx, domain.ID.String(), true)
		w.Write([]byte(`{"status": "verified"}`))
	} else {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(fmt.Sprintf(`{"status": "failed", "expected": "cortex-verification=%s"}`, domain.VerificationToken)))
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	// Simplified stats for Trending Cards
	// In a real app, this would query aggregated findings over time
	stats := map[string]interface{}{
		"total_assets":    12,
		"critical_risks":  3,
		"high_risks":      8,
		"scans_completed": 45,
		"trending_up":     true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
