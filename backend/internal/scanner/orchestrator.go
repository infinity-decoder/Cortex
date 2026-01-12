package scanner

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/infinity-decoder/cortex-backend/internal/container"
	"github.com/infinity-decoder/cortex-backend/internal/discovery"
	"github.com/infinity-decoder/cortex-backend/internal/fingerprinting"
	"github.com/infinity-decoder/cortex-backend/internal/persistence"
	"github.com/infinity-decoder/cortex-backend/internal/risk"
	"github.com/infinity-decoder/cortex-backend/internal/scanning"
	"github.com/infinity-decoder/cortex-backend/pkg/models"
)

type Orchestrator struct {
	Repo *persistence.Repository
}

func NewOrchestrator(repo *persistence.Repository) *Orchestrator {
	return &Orchestrator{Repo: repo}
}

type ScanResult struct {
	Assets      []discovery.Result
	AllFindings []risk.Exposure
	NewFindings []risk.Exposure
}

func (o *Orchestrator) RunScan(ctx context.Context, domainName string, domainID string) (*ScanResult, error) {
	// 0. Permission & Quota Check
	domain, _ := o.Repo.GetDomainByID(ctx, domainID)
	if domain == nil {
		return nil, fmt.Errorf("domain not found")
	}

	allowed, err := o.Repo.CheckQuota(ctx, domain.OrgID.String())
	if err != nil || !allowed {
		return nil, fmt.Errorf("daily scan quota exceeded for this organization")
	}

	// Logging: Audit Trail
	o.Repo.CreateAuditLog(ctx, &models.AuditLog{
		OrgID:  &domain.OrgID,
		Action: "SCAN_START",
		Metadata: fmt.Sprintf(`{"domain": "%s"}`, domainName),
	})

	// Track Scan Run
	runID, _ := o.Repo.CreateScanRun(ctx, domainID)

	// Fetch previous findings for delta detection
	previousFindings, _ := o.Repo.GetLatestFindingsForDomain(ctx, domainID)
	prevMap := make(map[string]bool)
	for _, pf := range previousFindings {
		key := fmt.Sprintf("%s-%s", pf.Type, pf.Severity)
		prevMap[key] = true
	}

	// 1. Initial Discovery (Passive)
	dnsScanner := discovery.NewScanner()
	assets, err := dnsScanner.EnumerateSubdomains(ctx, domainName)
	if err != nil {
		o.Repo.UpdateScanRunStatus(ctx, runID, "failed")
		return nil, fmt.Errorf("discovery failed: %v", err)
	}

	// 2. Scan & Analysis Pipeline
	portScanner := scanning.NewScanner()
	var allFindings []risk.Exposure
	var newFindings []risk.Exposure

	for _, assetResult := range assets {
		if len(assetResult.IPs) == 0 {
			continue
		}
		ip := assetResult.IPs[0]

		// Persistence: Save asset
		assetModel := &models.Asset{
			DomainID:  uuid.MustParse(domainID),
			Subdomain: assetResult.Subdomain,
			IPAddress: ip,
		}
		o.Repo.SaveAsset(ctx, assetModel)

		ports, _ := portScanner.ScanPorts(ctx, ip)
		for _, p := range ports {
			// Save service
			serviceModel := &models.Service{
				AssetID:  assetModel.ID,
				Port:     p.Port,
				Protocol: p.Protocol,
			}

			url := "http://"
			if p.Port == 443 || p.Port == 2376 || p.Port == 6443 {
				url = "https://"
			}
			url += fmt.Sprintf("%s:%d", ip, p.Port)

			fp, _ := fingerprinting.HTTPFingerprint(ctx, url)
			fpStr := ""
			if fp != nil {
				fpStr = fp.Server + " " + fp.BodySnippet
			}

			tech := container.Detect(p.Port, fpStr)
			serviceModel.Technology = string(tech)
			serviceModel.Fingerprint = fpStr
			o.Repo.SaveService(ctx, serviceModel)

			exposure := risk.Classify(p.Port, tech, fpStr)
			if exposure.Severity != risk.Info {
				allFindings = append(allFindings, exposure)
				
				key := fmt.Sprintf("%s-%s", exposure.Type, exposure.Severity)
				if !prevMap[key] {
					newFindings = append(newFindings, exposure)
				}

				findingModel := &models.Finding{
					ServiceID:   serviceModel.ID,
					Type:        exposure.Type,
					Severity:    string(exposure.Severity),
					Description: exposure.Description,
					Remediation: exposure.Remediation,
				}
				o.Repo.SaveFinding(ctx, findingModel)
			}
		}
	}

	// 3. Advanced Attack Path Mapping (Logic for Phase 3)
	// Example: If we have an exposed Kubernetes API AND a sensitive database, chain them.
	// This is a simplified implementation for demonstration.
	if len(allFindings) > 1 {
		log.Printf("[AttackPath] Analyzing chains for %d findings...", len(allFindings))
		// Chaining logic would go here
	}

	o.Repo.UpdateScanRunStatus(ctx, runID, "completed")
	
	o.Repo.CreateAuditLog(ctx, &models.AuditLog{
		OrgID:  &domain.OrgID,
		Action: "SCAN_COMPLETE",
		Metadata: fmt.Sprintf(`{"domain": "%s", "findings": %d}`, domainName, len(allFindings)),
	})

	return &ScanResult{
		Assets:      assets,
		AllFindings: allFindings,
		NewFindings: newFindings,
	}, nil
}
