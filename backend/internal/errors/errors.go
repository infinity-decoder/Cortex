package errors

import (
	"encoding/json"
	"net/http"
)

// APIError represents a standardized error response
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteError writes a standardized JSON error response
func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Code:    code,
		Message: message,
	})
}

// WriteErrorWithDetails writes a standardized JSON error response with additional details
func WriteErrorWithDetails(w http.ResponseWriter, status int, code, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Code:    code,
		Message: message,
		Details: details,
	})
}

// Common error codes
const (
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeValidationError   = "VALIDATION_ERROR"
	ErrCodeInternalError     = "INTERNAL_ERROR"
	ErrCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrCodeDomainNotVerified = "DOMAIN_NOT_VERIFIED"
	ErrCodeQuotaExceeded     = "QUOTA_EXCEEDED"
)
