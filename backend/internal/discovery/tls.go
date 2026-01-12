package discovery

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// ExtractSANs connects to a domain via TLS and retrieves Subject Alternative Names
func ExtractSANs(domain string) ([]string, error) {
	config := &tls.Config{
		InsecureSkipVerify: true, // We only want the cert info, no exploitation
	}

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:443", domain), config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sans []string
	state := conn.ConnectionState()
	for _, cert := range state.PeerCertificates {
		sans = append(sans, cert.DNSNames...)
	}

	// Remove duplicates
	uniqueSANs := make(map[string]bool)
	var result []string
	for _, name := range sans {
		if !uniqueSANs[name] {
			uniqueSANs[name] = true
			result = append(result, name)
		}
	}

	return result, nil
}
