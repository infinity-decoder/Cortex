package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infinity-decoder/cortex-backend/internal/container"
	"github.com/infinity-decoder/cortex-backend/internal/discovery"
	"github.com/infinity-decoder/cortex-backend/internal/fingerprinting"
	"github.com/infinity-decoder/cortex-backend/internal/risk"
	"github.com/infinity-decoder/cortex-backend/internal/scanning"
)

type ScanRequest struct {
	Domain string `json:"domain"`
}

type ScanResponse struct {
	Domain  string              `json:"domain"`
	Assets  []discovery.Result  `json:"assets"`
	Findings []risk.Exposure    `json:"findings"`
}

func handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	
	// 1. Discovery
	dnsScanner := discovery.NewScanner()
	assets, err := dnsScanner.EnumerateSubdomains(ctx, req.Domain)
	if err != nil {
		http.Error(w, "Discovery failed", http.StatusInternalServerError)
		return
	}

	// 2. Scanning & Fingerprinting
	portScanner := scanning.NewScanner()
	var allFindings []risk.Exposure

	for _, asset := range assets {
		// Use the first IP for demo purposes
		if len(asset.IPs) == 0 {
			continue
		}
		ip := asset.IPs[0]

		ports, _ := portScanner.ScanPorts(ctx, ip)
		for _, p := range ports {
			// Basic fingerprinting
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
			exposure := risk.Classify(p.Port, tech, fpStr)
			
			if exposure.Severity != risk.Info {
				allFindings = append(allFindings, exposure)
			}
		}
	}

	resp := ScanResponse{
		Domain:   req.Domain,
		Assets:   assets,
		Findings: allFindings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
