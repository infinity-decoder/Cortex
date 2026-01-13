package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
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
// PassiveDiscovery finds subdomains via crt.sh (Certificate Transparency logs)
func (s *Scanner) PassiveDiscovery(ctx context.Context, rootDomain string) ([]Result, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", rootDomain)
	
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var crtResults []struct {
		NameValue string `json:"name_value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&crtResults); err != nil {
		return nil, err
	}

	uniqueSubs := make(map[string]bool)
	for _, res := range crtResults {
		// crt.sh often returns wildcards or multiple names per entry
		names := strings.Split(res.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name == "" || strings.Contains(name, "*") {
				continue
			}
			// Only keep subdomains of the root
			if strings.HasSuffix(name, rootDomain) && name != rootDomain {
				sub := strings.TrimSuffix(name, "."+rootDomain)
				uniqueSubs[sub] = true
			}
		}
	}

	var results []Result
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Low concurrency for DNS resolution

	for sub := range uniqueSubs {
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

	wg.Wait()
	return results, nil
}
