package container

import (
	"strings"
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
