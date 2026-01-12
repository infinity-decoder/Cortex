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
		if err := o.Repo.SaveAsset(ctx, assetModel); err != nil {
			log.Printf("Warning: failed to save asset: %v", err)
		}

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
				
				// Delta Detection
				key := fmt.Sprintf("%s-%s", exposure.Type, exposure.Severity)
				if !prevMap[key] {
					newFindings = append(newFindings, exposure)
				}

				// Persistence: Save Finding
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

	return &ScanResult{
		Assets:      assets,
		AllFindings: allFindings,
		NewFindings: newFindings,
	}, nil
}
