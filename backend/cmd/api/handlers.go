package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cortex-backend/internal/discovery"
	"cortex-backend/internal/auth"
	"cortex-backend/internal/persistence"
	"cortex-backend/internal/risk"
	"cortex-backend/internal/scanner"
	"cortex-backend/internal/validation"
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
	
	// Validate domain
	if err := validation.ValidateDomain(req.Domain); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 1. Check if Domain is verified and belongs to Org
	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err != nil || !domain.Verified {
		// Auto-initialize if not exists, but block scan
		if err != nil {
			s.Repo.GetOrCreateDomain(ctx, orgID.String(), req.Domain)
		}
		
		http.Error(w, "Domain not verified. please verify ownership via DNS TXT record first.", http.StatusForbidden)
		return
	}

	// Audit log: Scan triggered
	ctx = auth.WithRequest(ctx, r)
	auth.LogAction(ctx, s.Repo, "SCAN_TRIGGERED", map[string]interface{}{
		"domain": req.Domain,
		"domain_id": domain.ID.String(),
	})

	// 2. Run Scan via Orchestrator
	result, err := s.Orchestrator.RunScan(ctx, req.Domain, domain.ID.String())
	if err != nil {
		// Audit log: Scan failed
		auth.LogAction(ctx, s.Repo, "SCAN_FAILED", map[string]interface{}{
			"domain": req.Domain,
			"domain_id": domain.ID.String(),
			"error": err.Error(),
		})
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Audit log: Scan completed
	auth.LogAction(ctx, s.Repo, "SCAN_COMPLETED", map[string]interface{}{
		"domain": req.Domain,
		"domain_id": domain.ID.String(),
		"assets_found": len(result.Assets),
		"findings_count": len(result.AllFindings),
	})

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
	
	// Validate domain
	if err := validation.ValidateDomain(req.Domain); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err != nil {
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	success, err := discovery.VerifyDomain(req.Domain, domain.VerificationToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification check failed: %v", err), http.StatusInternalServerError)
		return
	}

	ctx = auth.WithRequest(ctx, r)
	
	if success {
		s.Repo.UpdateDomainVerification(ctx, domain.ID.String(), true)
		
		// Audit log: Domain verification
		auth.LogAction(ctx, s.Repo, "DOMAIN_VERIFIED", map[string]interface{}{
			"domain": req.Domain,
			"domain_id": domain.ID.String(),
		})
		
		w.Write([]byte(`{"status": "verified"}`))
	} else {
		// Audit log: Domain verification failed
		auth.LogAction(ctx, s.Repo, "DOMAIN_VERIFY_FAILED", map[string]interface{}{
			"domain": req.Domain,
			"domain_id": domain.ID.String(),
		})
		
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
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
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
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
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
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
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

func (s *Server) handleUpdatePlan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := s.Repo.UpdateOrgPlan(ctx, orgID.String(), req.Plan)
	if err != nil {
		http.Error(w, "Failed to update plan", http.StatusInternalServerError)
		return
	}

	// Audit log: Plan update
	ctx = auth.WithRequest(ctx, r)
	auth.LogAction(ctx, s.Repo, "PLAN_UPDATED", map[string]interface{}{
		"org_id": orgID.String(),
		"new_plan": req.Plan,
	})

	w.Write([]byte(`{"status": "success"}`))
}

func (s *Server) handleGetPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	plan, err := s.Repo.GetOrgPlan(ctx, orgID.String())
	if err != nil {
		http.Error(w, "Failed to fetch plan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"plan": plan})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgName  string `json:"org_name"`
		FullName string `json:"full_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Validate inputs
	if err := validation.ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := validation.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := validation.ValidateOrgName(req.OrgName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := validation.ValidateFullName(req.FullName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ctx := auth.WithRequest(r.Context(), r)
	
	userID, err := s.Repo.CreateUser(ctx, req.Email, string(hashedPassword), req.FullName)
	if err != nil {
		http.Error(w, "user already exists or database error", http.StatusConflict)
		return
	}

	orgID, err := s.Repo.CreateOrganization(ctx, req.OrgName)
	if err != nil {
		http.Error(w, "failed to create organization", http.StatusInternalServerError)
		return
	}

	err = s.Repo.AddUserToOrganization(ctx, orgID, userID, "owner")
	if err != nil {
		http.Error(w, "failed to link user to organization", http.StatusInternalServerError)
		return
	}

	uUUID, _ := uuid.Parse(userID)
	oUUID, _ := uuid.Parse(orgID)
	token, _ := auth.GenerateToken(uUUID, oUUID)

	// Audit log: User registration
	auth.LogAction(ctx, s.Repo, "USER_REGISTER", map[string]interface{}{
		"user_id": userID,
		"email":    req.Email,
		"org_id":   orgID,
		"org_name": req.OrgName,
	})

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
	
	// Validate email format
	if err := validation.ValidateEmail(req.Email); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Basic password validation (not empty)
	if req.Password == "" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	user, err := s.Repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Don't reveal if user exists or not
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if account is locked
	ctx := r.Context()
	if err := auth.CheckLockout(ctx, s.Repo, user.ID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Record failed attempt
		auth.RecordFailedAttempt(ctx, s.Repo, user.ID)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Clear failed attempts on successful login
	auth.ClearFailedAttempts(ctx, s.Repo, user.ID)

	orgID, err := s.Repo.GetUserOrganization(ctx, user.ID.String())
	if err != nil {
		http.Error(w, "no organization found for user", http.StatusForbidden)
		return
	}

	oUUID, _ := uuid.Parse(orgID)
	token, _ := auth.GenerateToken(user.ID, oUUID)

	// Audit log: Successful login
	ctx = auth.WithRequest(ctx, r)
	auth.LogAction(ctx, s.Repo, "USER_LOGIN_SUCCESS", map[string]interface{}{
		"user_id": user.ID.String(),
		"email":   req.Email,
		"org_id":  orgID,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"status": "success",
	})
}

func (s *Server) handleGetDomains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domains, err := s.Repo.GetVerifiedDomains(ctx, orgID.String())
	if err != nil {
		http.Error(w, "Failed to fetch domains", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}
