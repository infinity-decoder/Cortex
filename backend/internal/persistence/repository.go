package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/infinity-decoder/cortex-backend/pkg/db"
	"github.com/infinity-decoder/cortex-backend/pkg/models"
)

type Repository struct {
	DB *db.Database
}

func NewRepository(database *db.Database) *Repository {
	return &Repository{DB: database}
}

// CreateScanRun initializes a new scan record
func (r *Repository) CreateScanRun(ctx context.Context, domainID string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO scan_runs (id, domain_id, status) VALUES ($1, $2, 'running')`
	_, err := r.DB.Pool.Exec(ctx, query, id, domainID)
	return id, err
}

// UpdateScanRunStatus updates the final status of a scan run
func (r *Repository) UpdateScanRunStatus(ctx context.Context, runID string, status string) error {
	query := `UPDATE scan_runs SET status = $1, finished_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.DB.Pool.Exec(ctx, query, status, runID)
	return err
}

// GetOrCreateDomain fetches a domain by name or creates it if it doesn't exist
func (r *Repository) GetOrCreateDomain(ctx context.Context, orgID string, domainName string) (string, error) {
	var id string
	checkQuery := `SELECT id FROM domains WHERE root_domain = $1 LIMIT 1`
	err := r.DB.Pool.QueryRow(ctx, checkQuery, domainName).Scan(&id)
	if err == nil {
		return id, nil
	}

	insertQuery := `INSERT INTO domains (org_id, root_domain, verification_token) VALUES ($1, $2, $3) RETURNING id`
	err = r.DB.Pool.QueryRow(ctx, insertQuery, orgID, domainName, uuid.New().String()).Scan(&id)
	return id, err
}

// SaveAsset saves or updates an asset (domain/IP)
func (r *Repository) SaveAsset(ctx context.Context, asset *models.Asset) error {
	query := `
		INSERT INTO assets (id, domain_id, subdomain, ip_address) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (domain_id, subdomain, ip_address) 
		DO UPDATE SET last_seen = CURRENT_TIMESTAMP 
		RETURNING id`
	if asset.ID == uuid.Nil {
		asset.ID = uuid.New()
	}
	err := r.DB.Pool.QueryRow(ctx, query, asset.ID, asset.DomainID, asset.Subdomain, asset.IPAddress).Scan(&asset.ID)
	return err
}

// SaveService saves or updates a discovered service
func (r *Repository) SaveService(ctx context.Context, service *models.Service) error {
	query := `
		INSERT INTO services (id, asset_id, port, protocol, fingerprint, technology) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		ON CONFLICT (asset_id, port, protocol) 
		DO UPDATE SET last_seen = CURRENT_TIMESTAMP, fingerprint = $5, technology = $6 
		RETURNING id`
	if service.ID == uuid.Nil {
		service.ID = uuid.New()
	}
	err := r.DB.Pool.QueryRow(ctx, query, service.ID, service.AssetID, service.Port, service.Protocol, service.Fingerprint, service.Technology).Scan(&service.ID)
	return err
}

// SaveFinding saves a discovered risk
func (r *Repository) SaveFinding(ctx context.Context, finding *models.Finding) error {
	query := `
		INSERT INTO findings (id, service_id, type, severity, description, remediation) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`
	if finding.ID == uuid.Nil {
		finding.ID = uuid.New()
	}
	err := r.DB.Pool.QueryRow(ctx, query, finding.ID, finding.ServiceID, finding.Type, finding.Severity, finding.Description, finding.Remediation).Scan(&finding.ID)
	return err
}

// GetLatestFindingsForDomain retrieves findings from the latest successful scan run
func (r *Repository) GetLatestFindingsForDomain(ctx context.Context, domainID string) ([]models.Finding, error) {
	query := `
		SELECT f.id, f.service_id, f.type, f.severity, f.description, f.remediation, f.first_seen, f.last_seen
		FROM findings f
		JOIN services s ON f.service_id = s.id
		JOIN assets a ON s.asset_id = a.id
		WHERE a.domain_id = $1
		ORDER BY f.last_seen DESC
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

// GetVerifiedDomains returns all domains that have been successfully verified for a specific org
func (r *Repository) GetVerifiedDomains(ctx context.Context, orgID string) ([]models.Domain, error) {
	query := `SELECT id, org_id, root_domain, verified, verification_token, created_at FROM domains WHERE verified = true AND org_id = $1`
	rows, err := r.DB.Pool.Query(ctx, query, orgID)
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

// GetDomainByNameAndOrg fetches a domain record by its root domain name and org ID
func (r *Repository) GetDomainByNameAndOrg(ctx context.Context, domainName string, orgID string) (*models.Domain, error) {
	var d models.Domain
	query := `SELECT id, org_id, root_domain, verified, verification_token, created_at FROM domains WHERE root_domain = $1 AND org_id = $2 LIMIT 1`
	err := r.DB.Pool.QueryRow(ctx, query, domainName, orgID).Scan(&d.ID, &d.OrgID, &d.RootDomain, &d.Verified, &d.VerificationToken, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// CheckQuota verifies if an organization has exceeded its daily scan limit
func (r *Repository) CheckQuota(ctx context.Context, orgID string) (bool, error) {
	var limit int
	var count int

	err := r.DB.Pool.QueryRow(ctx, "SELECT daily_scan_limit FROM organizations WHERE id = $1", orgID).Scan(&limit)
	if err != nil {
		return false, err
	}

	err = r.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM scan_runs sr JOIN domains d ON sr.domain_id = d.id WHERE d.org_id = $1 AND sr.started_at >= CURRENT_DATE", orgID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count < limit, nil
}

// CreateAuditLog records a sensitive action
func (r *Repository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	query := `INSERT INTO audit_logs (user_id, org_id, action, metadata, ip_address) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Pool.Exec(ctx, query, log.UserID, log.OrgID, log.Action, log.Metadata, log.IPAddress)
	return err
}

// GetDomainByID fetches a domain by ID
func (r *Repository) GetDomainByID(ctx context.Context, domainID string) (*models.Domain, error) {
	var d models.Domain
	query := `SELECT id, org_id, root_domain, verified, verification_token, created_at FROM domains WHERE id = $1`
	err := r.DB.Pool.QueryRow(ctx, query, domainID).Scan(&d.ID, &d.OrgID, &d.RootDomain, &d.Verified, &d.VerificationToken, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
// GetAssetsByDomain retrieves all discovered assets for a root domain
func (r *Repository) GetAssetsByDomain(ctx context.Context, domainID string) ([]models.Asset, error) {
	query := `SELECT id, domain_id, subdomain, ip_address, last_seen FROM assets WHERE domain_id = $1 ORDER BY last_seen DESC`
	rows, err := r.DB.Pool.Query(ctx, query, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var a models.Asset
		err := rows.Scan(&a.ID, &a.DomainID, &a.Subdomain, &a.IPAddress, &a.LastSeen)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// GetGlobalStats calculates real counts across the platform or org
func (r *Repository) GetGlobalStats(ctx context.Context, orgID string) (map[string]interface{}, error) {
	var totalAssets, criticalRisks, highRisks, scansCompleted int

	// Count Assets
	err := r.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM assets a JOIN domains d ON a.domain_id = d.id WHERE d.org_id = $1", orgID).Scan(&totalAssets)
	if err != nil { return nil, err }

	// Count Critical Risks
	err = r.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM findings f JOIN services s ON f.service_id = s.id JOIN assets a ON s.asset_id = a.id JOIN domains d ON a.domain_id = d.id WHERE d.org_id = $1 AND f.severity = 'critical'", orgID).Scan(&criticalRisks)
	if err != nil { return nil, err }

	// Count High Risks
	err = r.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM findings f JOIN services s ON f.service_id = s.id JOIN assets a ON s.asset_id = a.id JOIN domains d ON a.domain_id = d.id WHERE d.org_id = $1 AND f.severity = 'high'", orgID).Scan(&highRisks)
	if err != nil { return nil, err }

	// Count Scans
	err = r.DB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM scan_runs sr JOIN domains d ON sr.domain_id = d.id WHERE d.org_id = $1", orgID).Scan(&scansCompleted)
	if err != nil { return nil, err }

	return map[string]interface{}{
		"total_assets":    totalAssets,
		"critical_risks":  criticalRisks,
		"high_risks":      highRisks,
		"scans_completed": scansCompleted,
		"trending_up":     true,
	}, nil
}
// GetServicesByDomain retrieves all services for all assets of a domain
func (r *Repository) GetServicesByDomain(ctx context.Context, domainID string) ([]map[string]interface{}, error) {
	query := `
		SELECT s.port, s.protocol, s.fingerprint, s.technology, a.subdomain as asset_name, f.severity as risk
		FROM services s
		JOIN assets a ON s.asset_id = a.id
		LEFT JOIN findings f ON f.service_id = s.id
		WHERE a.domain_id = $1
		ORDER BY s.port ASC
	`
	rows, err := r.DB.Pool.Query(ctx, query, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []map[string]interface{}
	for rows.Next() {
		var port int
		var protocol, fingerprint, technology, assetName string
		var risk *string
		err := rows.Scan(&port, &protocol, &fingerprint, &technology, &assetName, &risk)
		if err != nil {
			return nil, err
		}
		
		rVal := "low"
		if risk != nil {
			rVal = *risk
		}

		services = append(services, map[string]interface{}{
			"port":       port,
			"proto":      protocol,
			"asset":      assetName,
			"service":    technology, // technology usually holds the service name
			"risk":       rVal,
			"fingerprint": fingerprint,
		})
	}
	return services, nil
}
// CreateUser registers a new user
func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, fullName string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Pool.Exec(ctx, query, id, email, passwordHash, fullName)
	return id, err
}

// GetUserByEmail finds a user by their email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	var u models.User
	var idStr string
	err := r.DB.Pool.QueryRow(ctx, query, email).Scan(&idStr, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	u.ID, _ = uuid.Parse(idStr)
	return &u, nil
}

// CreateOrganization creates a new org and returns its ID
func (r *Repository) CreateOrganization(ctx context.Context, name string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO organizations (id, name, plan_tier) VALUES ($1, $2, 'free')`
	_, err := r.DB.Pool.Exec(ctx, query, id, name)
	return id, err
}

// AddUserToOrganization links a user to an org with a role
func (r *Repository) AddUserToOrganization(ctx context.Context, orgID, userID, role string) error {
	query := `INSERT INTO org_users (org_id, user_id, role) VALUES ($1, $2, $3)`
	_, err := r.DB.Pool.Exec(ctx, query, orgID, userID, role)
	return err
}

// GetUserOrganization returns the first organization ID for a user
func (r *Repository) GetUserOrganization(ctx context.Context, userID string) (string, error) {
	var orgID string
	query := `SELECT org_id FROM org_users WHERE user_id = $1 LIMIT 1`
	err := r.DB.Pool.QueryRow(ctx, query, userID).Scan(&orgID)
	return orgID, err
}

// UpdateOrgPlan updates the plan tier for an organization
func (r *Repository) UpdateOrgPlan(ctx context.Context, orgID string, plan string) error {
	query := `UPDATE organizations SET plan_tier = $1 WHERE id = $2`
	_, err := r.DB.Pool.Exec(ctx, query, plan, orgID)
	return err
}

// GetOrgPlan retrieves the current plan tier for an organization
func (r *Repository) GetOrgPlan(ctx context.Context, orgID string) (string, error) {
	var plan string
	query := `SELECT plan_tier FROM organizations WHERE id = $1`
	err := r.DB.Pool.QueryRow(ctx, query, orgID).Scan(&plan)
	return plan, err
}
