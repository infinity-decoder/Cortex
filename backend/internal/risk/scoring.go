package risk

import (
	"github.com/infinity-decoder/cortex-backend/internal/container"
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
