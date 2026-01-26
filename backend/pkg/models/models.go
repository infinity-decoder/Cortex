package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type Organization struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	PlanTier string    `json:"planTier" db:"plan_tier"`
}

type Domain struct {
	ID                uuid.UUID `json:"id" db:"id"`
	OrgID             uuid.UUID `json:"orgId" db:"org_id"`
	RootDomain        string    `json:"rootDomain" db:"root_domain"`
	Verified          bool      `json:"verified" db:"verified"`
	VerificationToken string    `json:"verificationToken" db:"verification_token"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
}

type Asset struct {
	ID        uuid.UUID `json:"id" db:"id"`
	DomainID  uuid.UUID `json:"domainId" db:"domain_id"`
	Subdomain string    `json:"subdomain" db:"subdomain"`
	IPAddress string    `json:"ipAddress" db:"ip_address"`
	LastSeen  time.Time `json:"lastSeen" db:"last_seen"`
}

type Service struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AssetID     uuid.UUID `json:"assetId" db:"asset_id"`
	Port        int       `json:"port" db:"port"`
	Protocol    string    `json:"protocol" db:"protocol"`
	Fingerprint string    `json:"fingerprint" db:"fingerprint"`
	Technology  string    `json:"technology" db:"technology"`
}

type Finding struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ServiceID   uuid.UUID `json:"serviceId" db:"service_id"`
	Type        string    `json:"type" db:"type"`
	Severity    string    `json:"severity" db:"severity"`
	Description string    `json:"description" db:"description"`
	Remediation string    `json:"remediation" db:"remediation"`
	FirstSeen   time.Time `json:"firstSeen" db:"first_seen"`
	LastSeen    time.Time `json:"lastSeen" db:"last_seen"`
}

type ScanRun struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	DomainID   uuid.UUID  `json:"domainId" db:"domain_id"`
	Status     string     `json:"status" db:"status"`
	StartedAt  time.Time  `json:"startedAt" db:"started_at"`
	FinishedAt *time.Time `json:"finishedAt,omitempty" db:"finished_at"`
}

type AuditLog struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"userId"`
	OrgID     *uuid.UUID `json:"orgId"`
	Action    string     `json:"action"`
	Metadata  string     `json:"metadata"` // JSON string
	IPAddress string     `json:"ipAddress"`
	CreatedAt time.Time  `json:"createdAt"`
}
