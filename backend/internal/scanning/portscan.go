package scanning

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// Scanner handles port scanning for an IP/Asset
type Scanner struct {
	TargetPorts []int
	Timeout     time.Duration
}

func NewScanner() *Scanner {
	return &Scanner{
		TargetPorts: []int{
			80, 443,               // HTTP/HTTPS
			2375, 2376,           // Docker API
			6443, 8443,           // K8s API
			10250, 10255,         // Kubelet
			5000,                 // Registry
			3000, 8080, 9000, 9090, // Dashboards / Monitoring
		},
		Timeout: 2 * time.Second,
	}
}

// OpenPort represents a found open service
type OpenPort struct {
	Port     int
	Protocol string
}

// ScanPorts checks for open ports on an IP address
func (s *Scanner) ScanPorts(ctx context.Context, ip string) ([]OpenPort, error) {
	var results []OpenPort
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 20) // Moderate concurrency

	for _, port := range s.TargetPorts {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, s.Timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				results = append(results, OpenPort{
					Port:     port,
					Protocol: "tcp",
				})
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return results, nil
}
