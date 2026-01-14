package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"cortex-backend/internal/persistence"
)

const (
	// MaxFailedAttempts is the maximum number of failed login attempts before lockout
	MaxFailedAttempts = 5
	// LockoutDuration is how long an account stays locked
	LockoutDuration = 15 * time.Minute
)

// CheckLockout checks if an account is currently locked
func CheckLockout(ctx context.Context, repo *persistence.Repository, userID uuid.UUID) error {
	locked, lockedUntil, err := repo.GetAccountLockoutStatus(ctx, userID.String())
	if err != nil {
		return err
	}

	if locked {
		if lockedUntil != nil && time.Now().Before(*lockedUntil) {
			return errors.New("account is locked due to too many failed login attempts. Please try again later")
		}
		// Lockout expired, clear it
		repo.ClearFailedAttempts(ctx, userID.String())
	}

	return nil
}

// RecordFailedAttempt records a failed login attempt and locks account if threshold reached
func RecordFailedAttempt(ctx context.Context, repo *persistence.Repository, userID uuid.UUID) error {
	attempts, err := repo.IncrementFailedAttempts(ctx, userID.String())
	if err != nil {
		return err
	}

	// Lock account if threshold reached
	if attempts >= MaxFailedAttempts {
		lockedUntil := time.Now().Add(LockoutDuration)
		return repo.LockAccount(ctx, userID.String(), lockedUntil)
	}

	return nil
}

// ClearFailedAttempts clears failed login attempts (called on successful login)
func ClearFailedAttempts(ctx context.Context, repo *persistence.Repository, userID uuid.UUID) error {
	return repo.ClearFailedAttempts(ctx, userID.String())
}
