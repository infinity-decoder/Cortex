package scheduler

import (
	"context"
	"log"
	"time"

	"cortex-backend/internal/persistence"
	"cortex-backend/internal/scanner"
	"cortex-backend/pkg/models"
)

type Scheduler struct {
	Repo         *persistence.Repository
	Orchestrator *scanner.Orchestrator
	Interval     time.Duration
}

func NewScheduler(repo *persistence.Repository, orch *scanner.Orchestrator, interval time.Duration) *Scheduler {
	return &Scheduler{
		Repo:         repo,
		Orchestrator: orch,
		Interval:     interval,
	}
}

// Start initiates the background scanning loop
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	log.Printf("Scheduler started with interval: %v", s.Interval)

	// Run initial check immediately
	s.runPendingScans(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runPendingScans(ctx)
		}
	}
}

func (s *Scheduler) runPendingScans(ctx context.Context) {
	log.Println("Checking for pending scans for verified domains...")
	
	domains, err := s.Repo.GetAllVerifiedDomains(ctx)
	if err != nil {
		log.Printf("Scheduler error: failed to fetch verified domains: %v", err)
		return
	}

	for _, d := range domains {
		log.Printf("Triggering automated scan for: %s", d.RootDomain)
		go func(domain models.Domain) {
			_, scanErr := s.Orchestrator.RunScan(ctx, domain.RootDomain, domain.ID.String())
			if scanErr != nil {
				log.Printf("Automated scan failed for %s: %v", domain.RootDomain, scanErr)
			} else {
				log.Printf("Automated scan completed for %s", domain.RootDomain)
			}
		}(d)
	}
}
