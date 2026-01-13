package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infinity-decoder/cortex-backend/internal/discovery"
	"github.com/infinity-decoder/cortex-backend/internal/auth"
	"github.com/infinity-decoder/cortex-backend/internal/persistence"
	"github.com/infinity-decoder/cortex-backend/internal/risk"
	"github.com/infinity-decoder/cortex-backend/internal/scanner"
	"github.com/infinity-decoder/cortex-backend/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (s *Server) handleGetAssets(w http.ResponseWriter, r *http.Request) {
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		http.Error(w, "Domain query parameter required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domain, err := s.Repo.GetDomainByName(ctx, domainName)
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	assets, err := s.Repo.GetAssetsByDomain(ctx, domain.ID.String())
	if err != nil {
		http.Error(w, "Failed to fetch assets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func (s *Server) handleGetFindings(w http.ResponseWriter, r *http.Request) {
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		http.Error(w, "Domain query parameter required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domain, err := s.Repo.GetDomainByName(ctx, domainName)
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	findings, err := s.Repo.GetLatestFindingsForDomain(ctx, domain.ID.String())
	if err != nil {
		http.Error(w, "Failed to fetch findings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(findings)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "organization id not found in context", http.StatusUnauthorized)
		return
	}
	
	stats, err := s.Repo.GetGlobalStats(ctx, orgID.String())
	if err != nil {
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
func (s *Server) handleGetServices(w http.ResponseWriter, r *http.Request) {
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		http.Error(w, "Domain query parameter required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	domain, err := s.Repo.GetDomainByName(ctx, domainName)
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	services, err := s.Repo.GetServicesByDomain(ctx, domain.ID.String())
	if err != nil {
		http.Error(w, "Failed to fetch services", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgName  string `json:"org_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	userID, err := s.Repo.CreateUser(r.Context(), req.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "user already exists or database error", http.StatusConflict)
		return
	}

	orgID, err := s.Repo.CreateOrganization(r.Context(), req.OrgName)
	if err != nil {
		http.Error(w, "failed to create organization", http.StatusInternalServerError)
		return
	}

	err = s.Repo.AddUserToOrganization(r.Context(), orgID, userID, "owner")
	if err != nil {
		http.Error(w, "failed to link user to organization", http.StatusInternalServerError)
		return
	}

	uUUID, _ := uuid.Parse(userID)
	oUUID, _ := uuid.Parse(orgID)
	token, _ := auth.GenerateToken(uUUID, oUUID)

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"status": "success",
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := s.Repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	orgID, err := s.Repo.GetUserOrganization(r.Context(), user.ID.String())
	if err != nil {
		http.Error(w, "no organization found for user", http.StatusForbidden)
		return
	}

	oUUID, _ := uuid.Parse(orgID)
	token, _ := auth.GenerateToken(user.ID, oUUID)

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"status": "success",
	})
}
