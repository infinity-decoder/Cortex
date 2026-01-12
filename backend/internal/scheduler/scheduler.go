package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/infinity-decoder/cortex-backend/internal/persistence"
)

type Scheduler struct {
	Repo     *persistence.Repository
	Interval time.Duration
}

func NewScheduler(repo *persistence.Repository, interval time.Duration) *Scheduler {
	return &Scheduler{
		Repo:     repo,
		Interval: interval,
	}
}

// Start initiates the background scanning loop
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	log.Printf("Scheduler started with interval: %v", s.Interval)

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
	// For Phase 2, we'll fetch all verified domains and check if they need a scan
	// In a real SaaS, this would check the plan tier and last scan time
	// For now, we'll just log and placeholder the logic
	log.Println("Checking for pending scans...")
	
	// Example logic:
	// domains, _ := s.Repo.GetVerifiedDomains(ctx)
	// for _, d := range domains {
	//    if s.needsScan(d) {
	//        go s.triggerScan(ctx, d)
	//    }
	// }
}
