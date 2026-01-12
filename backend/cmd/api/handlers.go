package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/infinity-decoder/cortex-backend/internal/container"
	"github.com/infinity-decoder/cortex-backend/internal/discovery"
	"github.com/infinity-decoder/cortex-backend/internal/fingerprinting"
	"github.com/infinity-decoder/cortex-backend/internal/persistence"
	"github.com/infinity-decoder/cortex-backend/internal/risk"
	"github.com/infinity-decoder/cortex-backend/internal/scanning"
	"github.com/infinity-decoder/cortex-backend/pkg/models"
)

type Server struct {
	Repo *persistence.Repository
}

type ScanRequest struct {
	Domain string `json:"domain"`
}

type ScanResponse struct {
	Domain      string             `json:"domain"`
	Assets      []discovery.Result `json:"assets"`
	NewFindings []risk.Exposure    `json:"new_findings"`
	AllFindings []risk.Exposure    `json:"all_findings"`
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 1. Get/Create Domain record
	domainID, err := s.Repo.GetOrCreateDomain(ctx, "00000000-0000-0000-0000-000000000000", req.Domain)
	if err != nil {
		http.Error(w, "Domain initialization failed", http.StatusInternalServerError)
		return
	}

	// Fetch previous findings for delta detection
	previousFindings, _ := s.Repo.GetLatestFindingsForDomain(ctx, domainID)
	prevMap := make(map[string]bool)
	for _, pf := range previousFindings {
		key := fmt.Sprintf("%s-%s", pf.Type, pf.Severity)
		prevMap[key] = true
	}

	// 2. Initial Discovery (Passive)
	dnsScanner := discovery.NewScanner()
	assets, err := dnsScanner.EnumerateSubdomains(ctx, req.Domain)
	if err != nil {
		http.Error(w, "Discovery failed", http.StatusInternalServerError)
		return
	}

	// 3. Scan & Analysis Pipeline
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
		s.Repo.SaveAsset(ctx, assetModel)

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
			s.Repo.SaveService(ctx, serviceModel)

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
				s.Repo.SaveFinding(ctx, findingModel)
			}
		}
	}

	resp := ScanResponse{
		Domain:      req.Domain,
		Assets:      assets,
		NewFindings: newFindings,
		AllFindings: allFindings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
