package scanner

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	"cortex-backend/internal/alerting"
	"cortex-backend/internal/container"
	"cortex-backend/internal/discovery"
	"cortex-backend/internal/fingerprinting"
	"cortex-backend/internal/persistence"
	"cortex-backend/internal/risk"
	"cortex-backend/internal/scanning"
	"cortex-backend/pkg/models"
)

type Orchestrator struct {
	Repo         *persistence.Repository
	AlertHandler *alerting.AlertHandler
}

func NewOrchestrator(repo *persistence.Repository) *Orchestrator {
	return &Orchestrator{
		Repo:         repo,
		AlertHandler: alerting.NewAlertHandler(),
	}
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

	// 1. Discovery (Active + Passive + TLS Certificate Analysis)
	dnsScanner := discovery.NewScanner()
	activeAssets, _ := dnsScanner.EnumerateSubdomains(ctx, domainName)
	passiveAssets, _ := dnsScanner.PassiveDiscovery(ctx, domainName)
	
	// TLS Certificate Analysis - Extract SANs
	tlsSANs, err := discovery.ExtractSANs(domainName)
	if err == nil && len(tlsSANs) > 0 {
		log.Printf("[Discovery] Found %d SANs from TLS certificate", len(tlsSANs))
		// Convert SANs to discovery results by resolving to IPs
		for _, san := range tlsSANs {
			// Skip if it's the root domain (already discovered)
			if san == domainName {
				continue
			}
			// Resolve SAN to IPs using net.LookupHost
			ips, err := net.LookupHost(san)
			if err == nil && len(ips) > 0 {
				// Extract subdomain from SAN
				subdomain := san
				if strings.HasSuffix(san, "."+domainName) {
					subdomain = strings.TrimSuffix(san, "."+domainName)
				}
				passiveAssets = append(passiveAssets, discovery.Result{
					Subdomain: subdomain,
					IPs:       ips,
				})
			}
		}
	}
	
	// Merge assets
	assetMap := make(map[string]discovery.Result)
	for _, a := range activeAssets { assetMap[a.Subdomain] = a }
	for _, a := range passiveAssets {
		if existing, ok := assetMap[a.Subdomain]; ok {
			// Merge IPs
			ipMap := make(map[string]bool)
			for _, ip := range existing.IPs { ipMap[ip] = true }
			for _, ip := range a.IPs { ipMap[ip] = true }
			ips := []string{}
			for ip := range ipMap { ips = append(ips, ip) }
			existing.IPs = ips
			assetMap[a.Subdomain] = existing
		} else {
			assetMap[a.Subdomain] = a
		}
	}

	var assets []discovery.Result
	for _, a := range assetMap { assets = append(assets, a) }

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

			// Basic Classification
			exposure := risk.Classify(p.Port, tech, fpStr)
			
			// Advanced Probing (Phase 3)
			advSev, advTitle, advDesc := container.ProbeAdvanced(ctx, ip, p.Port, tech)
			if advSev != "" {
				exposure = risk.Exposure{
					Type:        advTitle,
					Severity:    risk.Severity(advSev),
					Description: advDesc,
					Remediation: exposure.Remediation, // Fallback to basic remediation
					AssetIP:     ip,
					Port:        p.Port,
					Technology:  string(tech),
				}
				// Specific remediation for advanced probes
				if advTitle == "Kubernetes Kubelet API Anonymous Access" {
					exposure.Remediation = "Set --anonymous-auth=false and --authorization-mode=Webhook in Kubelet configuration."
				} else if advTitle == "Exposed Docker Remote API (Unauthenticated)" {
					exposure.Remediation = "Disable TCP access to the Docker API or enforce MTLS authentication using certificates."
				}
			} else {
				// Add asset context to basic exposure
				exposure.AssetIP = ip
				exposure.Port = p.Port
				exposure.Technology = string(tech)
			}

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

	// 3. Advanced Attack Path Mapping
	if len(allFindings) > 1 {
		log.Printf("[AttackPath] Analyzing chains for %d findings...", len(allFindings))
		attackPaths := risk.AnalyzeAttackPaths(allFindings)
		if len(attackPaths) > 0 {
			log.Printf("[AttackPath] Identified %d attack paths", len(attackPaths))
			for _, path := range attackPaths {
				log.Printf("[AttackPath] Path: %s (Risk: %s, Score: %d) - %s", 
					path.ID, path.CombinedRisk, path.Score, path.Description)
			}
			// Store attack paths (for future: add to database)
			// For now, we log them and could add them to scan result metadata
		}
	}

	o.Repo.UpdateScanRunStatus(ctx, runID, "completed")
	
	o.Repo.CreateAuditLog(ctx, &models.AuditLog{
		OrgID:  &domain.OrgID,
		Action: "SCAN_COMPLETE",
		Metadata: fmt.Sprintf(`{"domain": "%s", "findings": %d}`, domainName, len(allFindings)),
	})

	// Send alerts for new critical/high findings
	if len(newFindings) > 0 {
		if err := o.AlertHandler.SendAlert(domainName, newFindings); err != nil {
			log.Printf("[Alert] Failed to send alert: %v", err)
		}
	}

	return &ScanResult{
		Assets:      assets,
		AllFindings: allFindings,
		NewFindings: newFindings,
	}, nil
}
