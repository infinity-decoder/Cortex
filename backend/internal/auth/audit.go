package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"cortex-backend/internal/persistence"
	"cortex-backend/pkg/models"
)

// LogAction logs an audit action
func LogAction(ctx context.Context, repo *persistence.Repository, action string, metadata map[string]interface{}) error {
	var userID *uuid.UUID
	var orgID *uuid.UUID

	// Extract user ID from context
	if uid, ok := ctx.Value(UserIDKey).(uuid.UUID); ok {
		userID = &uid
	}

	// Extract org ID from context
	if oid, ok := ctx.Value(OrgIDKey).(uuid.UUID); ok {
		orgID = &oid
	}

	// Extract IP address from request if available
	var ipAddress string
	if r, ok := ctx.Value("http_request").(*http.Request); ok {
		ipAddress = GetClientIP(r)
	}

	// Marshal metadata to JSON
	metadataJSON := "{}"
	if len(metadata) > 0 {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	log := &models.AuditLog{
		UserID:    userID,
		OrgID:     orgID,
		Action:    action,
		Metadata:  metadataJSON,
		IPAddress: ipAddress,
	}

	return repo.CreateAuditLog(ctx, log)
}

// GetClientIP extracts the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP (client IP)
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// WithRequest adds the HTTP request to context for audit logging
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, "http_request", r)
}
