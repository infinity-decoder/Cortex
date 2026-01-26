package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cortex-backend/internal/discovery"
	"cortex-backend/internal/auth"
	"cortex-backend/internal/errors"
	"cortex-backend/internal/persistence"
	"cortex-backend/internal/queue"
	"cortex-backend/internal/risk"
	"cortex-backend/internal/scanner"
	"cortex-backend/internal/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	Repo         *persistence.Repository
	Orchestrator *scanner.Orchestrator
	Queue        *queue.Queue
}

type ScanRequest struct {
	Domain string `json:"domain"`
}

type ScanResponse struct {
	Domain      string             `json:"domain"`
	Assets      []discovery.Result `json:"assets"`
	NewFindings []risk.Exposure    `json:"newFindings"`
	AllFindings []risk.Exposure    `json:"allFindings"`
	Verified    bool               `json:"verified"`
}

type VerifyRequest struct {
	Domain string `json:"domain"`
}

type CreateDomainRequest struct {
	Domain string `json:"domain"`
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}
	
	// Validate domain
	if err := validation.ValidateDomain(req.Domain); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	// 1. Check if Domain is verified and belongs to Org
	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err != nil || !domain.Verified {
		// Auto-initialize if not exists, but block scan
		if err != nil {
			s.Repo.GetOrCreateDomain(ctx, orgID.String(), req.Domain)
		}
		
		errors.WriteError(w, http.StatusForbidden, errors.ErrCodeDomainNotVerified, "Domain not verified. Please verify ownership via DNS TXT record first.")
		return
	}

	// Audit log: Scan triggered
	ctx = auth.WithRequest(ctx, r)
	auth.LogAction(ctx, s.Repo, "SCAN_TRIGGERED", map[string]interface{}{
		"domain": req.Domain,
		"domain_id": domain.ID.String(),
	})

	// 2. Queue Scan Job (Async Processing)
	jobID := s.Queue.Enqueue(req.Domain, domain.ID.String())
	
	// Audit log: Scan queued
	auth.LogAction(ctx, s.Repo, "SCAN_QUEUED", map[string]interface{}{
		"domain": req.Domain,
		"domain_id": domain.ID.String(),
		"job_id": jobID,
	})

	// Return job ID immediately
	resp := map[string]interface{}{
		"jobId": jobID,
		"status": "queued",
		"message": "Scan has been queued and will be processed shortly",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}
	
	// Validate domain
	if err := validation.ValidateDomain(req.Domain); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusNotFound, errors.ErrCodeNotFound, "Domain not found")
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
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, "Domain query parameter required")
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusNotFound, errors.ErrCodeNotFound, "Domain not found")
		return
	}

	assets, err := s.Repo.GetAssetsByDomain(ctx, domain.ID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch assets")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func (s *Server) handleGetFindings(w http.ResponseWriter, r *http.Request) {
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, "Domain query parameter required")
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusNotFound, errors.ErrCodeNotFound, "Domain not found")
		return
	}

	findings, err := s.Repo.GetLatestFindingsForDomain(ctx, domain.ID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch findings")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(findings)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Organization ID not found in context")
		return
	}
	
	stats, err := s.Repo.GetGlobalStats(ctx, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch stats")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
func (s *Server) handleGetServices(w http.ResponseWriter, r *http.Request) {
	domainName := r.URL.Query().Get("domain")
	if domainName == "" {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, "Domain query parameter required")
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, domainName, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusNotFound, errors.ErrCodeNotFound, "Domain not found")
		return
	}

	services, err := s.Repo.GetServicesByDomain(ctx, domain.ID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch services")
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
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	err := s.Repo.UpdateOrgPlan(ctx, orgID.String(), req.Plan)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to update plan")
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
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	plan, err := s.Repo.GetOrgPlan(ctx, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch plan")
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
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}

	// Validate inputs
	if err := validation.ValidateEmail(req.Email); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}
	
	if err := validation.ValidatePassword(req.Password); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}
	
	if err := validation.ValidateOrgName(req.OrgName); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}
	
	if err := validation.ValidateFullName(req.FullName); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Internal server error")
		return
	}

	ctx := auth.WithRequest(r.Context(), r)
	
	userID, err := s.Repo.CreateUser(ctx, req.Email, string(hashedPassword), req.FullName)
	if err != nil {
		errors.WriteError(w, http.StatusConflict, "USER_EXISTS", "User already exists or database error")
		return
	}

	orgID, err := s.Repo.CreateOrganization(ctx, req.OrgName)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to create organization")
		return
	}

	err = s.Repo.AddUserToOrganization(ctx, orgID, userID, "owner")
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to link user to organization")
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
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}
	
	// Validate email format
	if err := validation.ValidateEmail(req.Email); err != nil {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Invalid credentials")
		return
	}
	
	// Basic password validation (not empty)
	if req.Password == "" {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Invalid credentials")
		return
	}

	user, err := s.Repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Don't reveal if user exists or not
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Invalid credentials")
		return
	}

	// Check if account is locked
	ctx := r.Context()
	if err := auth.CheckLockout(ctx, s.Repo, user.ID); err != nil {
		errors.WriteError(w, http.StatusForbidden, "ACCOUNT_LOCKED", err.Error())
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Record failed attempt
		auth.RecordFailedAttempt(ctx, s.Repo, user.ID)
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Invalid credentials")
		return
	}

	// Clear failed attempts on successful login
	auth.ClearFailedAttempts(ctx, s.Repo, user.ID)

	orgID, err := s.Repo.GetUserOrganization(ctx, user.ID.String())
	if err != nil {
		errors.WriteError(w, http.StatusForbidden, errors.ErrCodeForbidden, "No organization found for user")
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

func (s *Server) handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	var req CreateDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeInvalidRequest, "Invalid request body")
		return
	}

	// Validate domain
	if err := validation.ValidateDomain(req.Domain); err != nil {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, err.Error())
		return
	}

	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	// Check if domain already exists
	existingDomain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err == nil && existingDomain != nil {
		// Domain exists, return it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":                existingDomain.ID.String(),
			"rootDomain":        existingDomain.RootDomain,
			"verified":          existingDomain.Verified,
			"verificationToken": existingDomain.VerificationToken,
			"message":           "Domain already exists",
		})
		return
	}

	// Create new domain
	domainID, err := s.Repo.GetOrCreateDomain(ctx, orgID.String(), req.Domain)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to create domain")
		return
	}

	// Fetch the created domain to get verification token
	domain, err := s.Repo.GetDomainByNameAndOrg(ctx, req.Domain, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch created domain")
		return
	}

	// Audit log
	ctx = auth.WithRequest(ctx, r)
	auth.LogAction(ctx, s.Repo, "DOMAIN_CREATED", map[string]interface{}{
		"domain":     req.Domain,
		"domain_id":   domainID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                domain.ID.String(),
		"rootDomain":        domain.RootDomain,
		"verified":          domain.Verified,
		"verificationToken":  domain.VerificationToken,
		"message":           "Domain created successfully. Please verify ownership via DNS TXT record.",
	})
}

func (s *Server) handleGetDomains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domains, err := s.Repo.GetVerifiedDomains(ctx, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch domains")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

func (s *Server) handleGetScanStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobId")
	if jobID == "" {
		errors.WriteError(w, http.StatusBadRequest, errors.ErrCodeValidationError, "jobId query parameter required")
		return
	}

	job, exists := s.Queue.GetJob(jobID)
	if !exists {
		errors.WriteError(w, http.StatusNotFound, errors.ErrCodeNotFound, "Job not found")
		return
	}

	response := map[string]interface{}{
		"jobId": job.ID,
		"status": job.Status,
		"domain": job.Domain,
		"createdAt": job.CreatedAt,
	}

	if job.StartedAt != nil {
		response["startedAt"] = job.StartedAt
	}
	if job.FinishedAt != nil {
		response["finishedAt"] = job.FinishedAt
	}
	if job.Status == "completed" && job.Result != nil {
		response["result"] = job.Result
	}
	if job.Status == "failed" && job.Error != nil {
		response["error"] = job.Error.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetAllDomains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orgID, ok := ctx.Value(auth.OrgIDKey).(uuid.UUID)
	if !ok {
		errors.WriteError(w, http.StatusUnauthorized, errors.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	domains, err := s.Repo.GetAllDomainsByOrg(ctx, orgID.String())
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, errors.ErrCodeInternalError, "Failed to fetch domains")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}
