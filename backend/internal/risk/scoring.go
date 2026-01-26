package risk

import (
	"cortex-backend/internal/container"
)

type Severity string

const (
	Critical Severity = "critical"
	High     Severity = "high"
	Medium   Severity = "medium"
	Low      Severity = "low"
	Info     Severity = "info"
)

type Exposure struct {
	Type        string
	Severity    Severity
	Description string
	Remediation string
	AssetIP     string // IP address of the asset
	Port        int    // Port number
	Technology  string // Technology detected
}

// AttackPath represents a chain of vulnerabilities that could be exploited together
type AttackPath struct {
	ID          string
	Steps       []Exposure
	CombinedRisk Severity
	Description string
	Score       int // 0-100 risk score
}

// AnalyzeAttackPaths identifies potential attack paths by chaining related findings
func AnalyzeAttackPaths(findings []Exposure) []AttackPath {
	var paths []AttackPath
	
	if len(findings) < 2 {
		return paths // Need at least 2 findings for a path
	}

	// Group findings by asset (IP address)
	assetMap := make(map[string][]Exposure)
	for _, f := range findings {
		if f.AssetIP != "" {
			assetMap[f.AssetIP] = append(assetMap[f.AssetIP], f)
		}
	}

	// Identify common attack patterns
	// Pattern 1: Kubernetes API + Exposed Database
	kubernetesFindings := []Exposure{}
	databaseFindings := []Exposure{}
	dockerFindings := []Exposure{}
	
	for _, f := range findings {
		if f.Technology == "kubernetes" || contains(f.Type, "Kubernetes", "Kubelet", "K8s") {
			kubernetesFindings = append(kubernetesFindings, f)
		}
		if contains(f.Type, "Database", "PostgreSQL", "MySQL", "MongoDB", "Redis") {
			databaseFindings = append(databaseFindings, f)
		}
		if f.Technology == "docker" || contains(f.Type, "Docker") {
			dockerFindings = append(dockerFindings, f)
		}
	}

	// Create attack path: K8s API + Database
	if len(kubernetesFindings) > 0 && len(databaseFindings) > 0 {
		path := AttackPath{
			ID: "k8s-db-chain",
			Steps: []Exposure{
				kubernetesFindings[0],
				databaseFindings[0],
			},
			Description: "Kubernetes API exposure combined with database access could allow cluster compromise and data exfiltration",
			Score: calculatePathScore([]Severity{kubernetesFindings[0].Severity, databaseFindings[0].Severity}),
		}
		path.CombinedRisk = scoreToSeverity(path.Score)
		paths = append(paths, path)
	}

	// Create attack path: Docker API + Database
	if len(dockerFindings) > 0 && len(databaseFindings) > 0 {
		path := AttackPath{
			ID: "docker-db-chain",
			Steps: []Exposure{
				dockerFindings[0],
				databaseFindings[0],
			},
			Description: "Docker API exposure combined with database access could allow container escape and data breach",
			Score: calculatePathScore([]Severity{dockerFindings[0].Severity, databaseFindings[0].Severity}),
		}
		path.CombinedRisk = scoreToSeverity(path.Score)
		paths = append(paths, path)
	}

	// Pattern 2: Multiple critical findings on same asset
	for ip, assetFindings := range assetMap {
		if len(assetFindings) >= 2 {
			criticalCount := 0
			for _, f := range assetFindings {
				if f.Severity == Critical {
					criticalCount++
				}
			}
			if criticalCount >= 2 {
				path := AttackPath{
					ID: "multi-critical-" + ip,
					Steps: assetFindings,
					Description: "Multiple critical vulnerabilities on the same asset increase exploitation likelihood",
					Score: calculatePathScore(getSeverities(assetFindings)),
				}
				path.CombinedRisk = scoreToSeverity(path.Score)
				paths = append(paths, path)
			}
		}
	}

	return paths
}

// Helper functions
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func calculatePathScore(severities []Severity) int {
	score := 0
	for _, sev := range severities {
		switch sev {
		case Critical:
			score += 40
		case High:
			score += 25
		case Medium:
			score += 15
		case Low:
			score += 5
		}
	}
	// Bonus for multiple critical/high findings
	if len(severities) >= 2 {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return score
}

func scoreToSeverity(score int) Severity {
	if score >= 70 {
		return Critical
	} else if score >= 50 {
		return High
	} else if score >= 30 {
		return Medium
	} else if score >= 10 {
		return Low
	}
	return Info
}

func getSeverities(findings []Exposure) []Severity {
	severities := make([]Severity, len(findings))
	for i, f := range findings {
		severities[i] = f.Severity
	}
	return severities
}

// Classify assigns risk to a detected service
func Classify(port int, tech container.Technology, fingerprint string) Exposure {
	// Example Logic: Docker API exposed on 2375 (no TLS) is Critical
	if port == 2375 && tech == container.Docker {
		return Exposure{
			Type:        "Insecure Docker API",
			Severity:    Critical,
			Description: "Docker Remote API is exposed without TLS authentication. An attacker can gain full control over the host.",
			Remediation: "Disable the Remote API or enable TLS authentication and restrict access to specific IPs.",
		}
	}

	// Kubernetes API exposed
	if tech == container.Kubernetes {
		return Exposure{
			Type:        "Exposed Kubernetes API",
			Severity:    High,
			Description: "A Kubernetes API server was detected. If misconfigured, it could allow unauthorized access to the cluster.",
			Remediation: "Ensure the API server requires authentication and is not accessible from the public internet.",
		}
	}

	// Default
	return Exposure{
		Type:        "Exposed Service",
		Severity:    Info,
		Description: "A service was detected on an open port.",
		Remediation: "Verify if this service is intended to be public.",
	}
}
