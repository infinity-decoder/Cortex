package container

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Technology string

const (
	Docker     Technology = "docker"
	Kubernetes Technology = "kubernetes"
	Registry   Technology = "registry"
	Unknown    Technology = "unknown"
)

type Signature struct {
	Tech          Technology
	Ports         []int
	Fingerprints  []string // Substrings to look for in HTTP headers/body
}

var Signatures = []Signature{
	{
		Tech:  Docker,
		Ports: []int{2375, 2376},
		Fingerprints: []string{
			"Docker",
			"Api-Version",
		},
	},
	{
		Tech:  Kubernetes,
		Ports: []int{6443, 8443, 10250, 10255},
		Fingerprints: []string{
			"k8s",
			"kubernetes",
			"kube",
			"Unauthorized", // Often returned by K8s API when exposed
		},
	},
	{
		Tech:  Registry,
		Ports: []int{5000},
		Fingerprints: []string{
			"Docker-Distribution-Api-Version",
		},
	},
}

// Detect identifies container technology based on port and fingerprint
func Detect(port int, fingerprint string) Technology {
	for _, sig := range Signatures {
		// Match by port
		portMatch := false
		for _, p := range sig.Ports {
			if p == port {
				portMatch = true
				break
			}
		}

		// Match by fingerprint
		fpMatch := false
		for _, s := range sig.Fingerprints {
			if strings.Contains(strings.ToLower(fingerprint), strings.ToLower(s)) {
				fpMatch = true
				break
			}
		}

		if portMatch || fpMatch {
			return sig.Tech
		}
	}
	return Unknown
}
// ProbeAdvanced checks for specific vulnerabilities once a technology is detected
func ProbeAdvanced(ctx context.Context, host string, port int, tech Technology) (string, string, string) {
	if tech == Kubernetes {
		return ProbeK8s(ctx, host, port)
	}
	if tech == Docker {
		return ProbeDocker(ctx, host, port)
	}
	return "", "", ""
}

func ProbeK8s(ctx context.Context, host string, port int) (string, string, string) {
	// Check for anonymous access to /pods or /api
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	
	url := fmt.Sprintf("https://%s:%d/pods", host, port)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			return "critical", "Kubernetes Kubelet API Anonymous Access", 
				"The Kubelet API allows anonymous users to list pods. This can lead to sensitive information disclosure and potential execution of commands in pods."
		}
	}
	
	return "", "", ""
}

func ProbeDocker(ctx context.Context, host string, port int) (string, string, string) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://%s:%d/v1.24/version", host, port)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			return "critical", "Exposed Docker Remote API (Unauthenticated)", 
				"The Docker Remote API is accessible without authentication. Attackers can execute commands and pull/push images."
		}
	}
	
	return "", "", ""
}
