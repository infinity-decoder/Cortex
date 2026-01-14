package alerting

import (
	"fmt"
	"log"

	"cortex-backend/internal/risk"
)

// AlertHandler handles sending alerts via different channels
type AlertHandler struct {
	// Add SMTP/Webhook config here later
}

func NewAlertHandler() *AlertHandler {
	return &AlertHandler{}
}

// SendAlert sends a notification about new high-risk findings
func (a *AlertHandler) SendAlert(domain string, findings []risk.Exposure) error {
	if len(findings) == 0 {
		return nil
	}

	// For Phase 2 MVP: Log the alert to terminal (simulating email/webhook)
	log.Printf("🚨 ALERT: New findings for %s:", domain)
	for _, f := range findings {
		log.Printf("- [%s] %s: %s", f.Severity, f.Type, f.Description)
	}

	// Logic for SMTP would go here
	// return smtp.SendMail(...)

	return nil
}

// FormatEmailBody creates a simple HTML/Text body for the finding
func FormatEmailBody(finding risk.Exposure) string {
	return fmt.Sprintf("Type: %s\nSeverity: %s\nDescription: %s\nRemediation: %s",
		finding.Type, finding.Severity, finding.Description, finding.Remediation)
}
