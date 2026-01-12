package discovery

import (
	"fmt"
	"net"
	"strings"
)

// VerifyDomain checks if a specific verification token exists in the DNS TXT records for a domain.
// Standard SaaS practice: check for cortex-verification=TOKEN
func VerifyDomain(domain, token string) (bool, error) {
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		return false, fmt.Errorf("failed to lookup TXT records: %v", err)
	}

	expected := fmt.Sprintf("cortex-verification=%s", token)
	for _, txt := range txtrecords {
		if strings.TrimSpace(txt) == expected {
			return true, nil
		}
	}

	return false, nil
}
