package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Organization struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	PlanTier string    `json:"plan_tier" db:"plan_tier"`
}

type Domain struct {
	ID                uuid.UUID `json:"id" db:"id"`
	OrgID             uuid.UUID `json:"org_id" db:"org_id"`
	RootDomain        string    `json:"root_domain" db:"root_domain"`
	Verified          bool      `json:"verified" db:"verified"`
	VerificationToken string    `json:"verification_token" db:"verification_token"`
}

type Asset struct {
	ID        uuid.UUID `json:"id" db:"id"`
	DomainID  uuid.UUID `json:"domain_id" db:"domain_id"`
	Subdomain string    `json:"subdomain" db:"subdomain"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
}

type Service struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AssetID     uuid.UUID `json:"asset_id" db:"asset_id"`
	Port        int       `json:"port" db:"port"`
	Protocol    string    `json:"protocol" db:"protocol"`
	Fingerprint string    `json:"fingerprint" db:"fingerprint"`
	Technology  string    `json:"technology" db:"technology"`
}

type Finding struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ServiceID   uuid.UUID `json:"service_id" db:"service_id"`
	Type        string    `json:"type" db:"type"`
	Severity    string    `json:"severity" db:"severity"`
	Description string    `json:"description" db:"description"`
	Remediation string    `json:"remediation" db:"remediation"`
	FirstSeen   time.Time `json:"first_seen" db:"first_seen"`
	LastSeen    time.Time `json:"last_seen" db:"last_seen"`
}

type ScanRun struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	DomainID   uuid.UUID  `json:"domain_id" db:"domain_id"`
	Status     string     `json:"status" db:"status"`
	StartedAt  time.Time  `json:"started_at" db:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty" db:"finished_at"`
}
