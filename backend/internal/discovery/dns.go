package discovery

import (
	"context"
	"fmt"
	"net"
	"sync"
)

// Scanner handles asset discovery for a domain
type Scanner struct {
	CommonSubdomains []string
}

func NewScanner() *Scanner {
	return &Scanner{
		CommonSubdomains: []string{
			"www", "api", "dev", "staging", "prod", "test",
			"db", "database", "admin", "dashboard", "portal",
			"git", "gitlab", "jenkins", "docker", "k8s", "kube",
			"registry", "vault", "vpn", "mail", "remote",
		},
	}
}

// Result represents a discovered asset
type Result struct {
	Subdomain string
	IPs       []string
}

// EnumerateSubdomains performs basic DNS brute-forcing (passive/safe)
func (s *Scanner) EnumerateSubdomains(ctx context.Context, rootDomain string) ([]Result, error) {
	var results []Result
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrency to be polite and avoid rate limits
	semaphore := make(chan struct{}, 10)

	for _, sub := range s.CommonSubdomains {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fullDomain := fmt.Sprintf("%s.%s", sub, rootDomain)
			ips, err := net.LookupHost(fullDomain)
			if err == nil && len(ips) > 0 {
				mu.Lock()
				results = append(results, Result{
					Subdomain: sub,
					IPs:       ips,
				})
				mu.Unlock()
			}
		}(sub)
	}

	// Also check the root domain itself
	wg.Add(1)
	go func() {
		defer wg.Done()
		ips, err := net.LookupHost(rootDomain)
		if err == nil && len(ips) > 0 {
			mu.Lock()
			results = append(results, Result{
				Subdomain: "",
				IPs:       ips,
			})
			mu.Unlock()
		}
	}()

	wg.Wait()
	return results, nil
}
