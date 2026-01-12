package persistence

import (
	"context"

	"github.com/infinity-decoder/cortex-backend/pkg/db"
	"github.com/infinity-decoder/cortex-backend/pkg/models"
)

type Repository struct {
	DB *db.Database
}

func NewRepository(database *db.Database) *Repository {
	return &Repository{DB: database}
}

// SaveScanRun initializes a scan run in the database
func (r *Repository) SaveScanRun(ctx context.Context, run *models.ScanRun) error {
	query := `INSERT INTO scan_runs (domain_id, status) VALUES ($1, $2) RETURNING id, started_at`
	err := r.DB.Pool.QueryRow(ctx, query, run.DomainID, run.Status).Scan(&run.ID, &run.StartedAt)
	return err
}

// GetOrCreateDomain fetches a domain by name or creates it if it doesn't exist
func (r *Repository) GetOrCreateDomain(ctx context.Context, orgID string, domainName string) (string, error) {
	var id string
	
	// For MVP, we'll use a fixed ORG ID or create one if needed.
	// Let's just find by name for simplicity if unique.
	checkQuery := `SELECT id FROM domains WHERE root_domain = $1 LIMIT 1`
	err := r.DB.Pool.QueryRow(ctx, checkQuery, domainName).Scan(&id)
	if err == nil {
		return id, nil
	}

	insertQuery := `INSERT INTO domains (org_id, root_domain, verification_token) VALUES ($1, $2, $3) RETURNING id`
	err = r.DB.Pool.QueryRow(ctx, insertQuery, "00000000-0000-0000-0000-000000000000", domainName, "dummy-token").Scan(&id)
	return id, err
}

// UpdateScanRunStatus updates the final status of a scan run
func (r *Repository) UpdateScanRunStatus(ctx context.Context, runID string, status string) error {
	query := `UPDATE scan_runs SET status = $1, finished_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.DB.Pool.Exec(ctx, query, status, runID)
	return err
}

// SaveAsset saves or updates an asset (domain/IP)
func (r *Repository) SaveAsset(ctx context.Context, asset *models.Asset) error {
	query := `
		INSERT INTO assets (domain_id, subdomain, ip_address) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (domain_id, subdomain, ip_address) 
		DO UPDATE SET last_seen = CURRENT_TIMESTAMP 
		RETURNING id`
	err := r.DB.Pool.QueryRow(ctx, query, asset.DomainID, asset.Subdomain, asset.IPAddress).Scan(&asset.ID)
	return err
}

// SaveService saves or updates a discovered service
func (r *Repository) SaveService(ctx context.Context, service *models.Service) error {
	query := `
		INSERT INTO services (asset_id, port, protocol, fingerprint, technology) 
		VALUES ($1, $2, $3, $4, $5) 
		ON CONFLICT (asset_id, port, protocol) 
		DO UPDATE SET last_seen = CURRENT_TIMESTAMP, fingerprint = $4, technology = $5 
		RETURNING id`
	err := r.DB.Pool.QueryRow(ctx, query, service.AssetID, service.Port, service.Protocol, service.Fingerprint, service.Technology).Scan(&service.ID)
	return err
}

// SaveFinding saves a discovered risk
func (r *Repository) SaveFinding(ctx context.Context, finding *models.Finding) error {
	query := `
		INSERT INTO findings (service_id, type, severity, description, remediation) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`
	err := r.DB.Pool.QueryRow(ctx, query, finding.ServiceID, finding.Type, finding.Severity, finding.Description, finding.Remediation).Scan(&finding.ID)
	return err
}

// GetLatestFindingsForDomain retrieves findings from the last successful scan run for a domain
func (r *Repository) GetLatestFindingsForDomain(ctx context.Context, domainID string) ([]models.Finding, error) {
	query := `
		SELECT f.id, f.service_id, f.type, f.severity, f.description, f.remediation, f.first_seen, f.last_seen
		FROM findings f
		JOIN services s ON f.service_id = s.id
		JOIN assets a ON s.asset_id = a.id
		WHERE a.domain_id = $1
		AND f.last_seen = (
			SELECT MAX(last_seen) FROM findings f2
			JOIN services s2 ON f2.service_id = s2.id
			JOIN assets a2 ON s2.asset_id = a2.id
			WHERE a2.domain_id = $1
		)
	`
	rows, err := r.DB.Pool.Query(ctx, query, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var findings []models.Finding
	for rows.Next() {
		var f models.Finding
		err := rows.Scan(&f.ID, &f.ServiceID, &f.Type, &f.Severity, &f.Description, &f.Remediation, &f.FirstSeen, &f.LastSeen)
		if err != nil {
			return nil, err
		}
		findings = append(findings, f)
	}
	return findings, nil
}

// GetVerifiedDomains returns all domains that have been successfully verified
func (r *Repository) GetVerifiedDomains(ctx context.Context) ([]models.Domain, error) {
	query := `SELECT id, org_id, root_domain, verified, verification_token, created_at FROM domains WHERE verified = true`
	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var d models.Domain
		err := rows.Scan(&d.ID, &d.OrgID, &d.RootDomain, &d.Verified, &d.VerificationToken, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, nil
}

// UpdateDomainVerification updates the verification status of a domain
func (r *Repository) UpdateDomainVerification(ctx context.Context, domainID string, verified bool) error {
	query := `UPDATE domains SET verified = $1 WHERE id = $2`
	_, err := r.DB.Pool.Exec(ctx, query, verified, domainID)
	return err
}

// GetDomainByName fetches a domain record by its root domain name
func (r *Repository) GetDomainByName(ctx context.Context, domainName string) (*models.Domain, error) {
	var d models.Domain
	query := `SELECT id, org_id, root_domain, verified, verification_token, created_at FROM domains WHERE root_domain = $1 LIMIT 1`
	err := r.DB.Pool.QueryRow(ctx, query, domainName).Scan(&d.ID, &d.OrgID, &d.RootDomain, &d.Verified, &d.VerificationToken, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
