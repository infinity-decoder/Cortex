package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"cortex-backend/internal/risk"
)

// AlertHandler handles sending alerts via different channels
type AlertHandler struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	WebhookURL   string
}

// NewAlertHandler creates a new alert handler with configuration from environment
func NewAlertHandler() *AlertHandler {
	return &AlertHandler{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		WebhookURL:   os.Getenv("WEBHOOK_URL"),
	}
}

// SendAlert sends a notification about new high-risk findings
func (a *AlertHandler) SendAlert(domain string, findings []risk.Exposure) error {
	if len(findings) == 0 {
		return nil
	}

	// Filter only critical and high severity findings
	criticalFindings := []risk.Exposure{}
	for _, f := range findings {
		if f.Severity == risk.Critical || f.Severity == risk.High {
			criticalFindings = append(criticalFindings, f)
		}
	}

	if len(criticalFindings) == 0 {
		return nil
	}

	// Log alert
	log.Printf("🚨 ALERT: %d critical/high findings for %s", len(criticalFindings), domain)
	for _, f := range criticalFindings {
		log.Printf("- [%s] %s: %s", f.Severity, f.Type, f.Description)
	}

	// Send email if configured
	if a.SMTPHost != "" && a.SMTPUser != "" {
		if err := a.sendEmail(domain, criticalFindings); err != nil {
			log.Printf("Failed to send email alert: %v", err)
		}
	}

	// Send webhook if configured
	if a.WebhookURL != "" {
		if err := a.sendWebhook(domain, criticalFindings); err != nil {
			log.Printf("Failed to send webhook alert: %v", err)
		}
	}

	return nil
}

// sendEmail sends an email alert
func (a *AlertHandler) sendEmail(domain string, findings []risk.Exposure) error {
	// Get recipient from environment (for MVP, use SMTP_FROM or SMTP_TO)
	to := os.Getenv("SMTP_TO")
	if to == "" {
		to = a.SMTPFrom
	}

	subject := fmt.Sprintf("🚨 Security Alert: %d Critical Findings for %s", len(findings), domain)
	body := a.formatEmailBody(domain, findings)

	// SMTP configuration
	addr := fmt.Sprintf("%s:%s", a.SMTPHost, a.SMTPPort)
	auth := smtp.PlainAuth("", a.SMTPUser, a.SMTPPassword, a.SMTPHost)

	// Email message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		to, subject, body))

	// Send email
	err := smtp.SendMail(addr, auth, a.SMTPFrom, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email alert sent to %s", to)
	return nil
}

// sendWebhook sends a webhook alert
func (a *AlertHandler) sendWebhook(domain string, findings []risk.Exposure) error {
	payload := map[string]interface{}{
		"domain":   domain,
		"timestamp": time.Now().Format(time.RFC3339),
		"count":     len(findings),
		"findings":  findings,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req, err := http.NewRequest("POST", a.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	log.Printf("Webhook alert sent to %s", a.WebhookURL)
	return nil
}

// formatEmailBody creates an HTML email body for the findings
func (a *AlertHandler) formatEmailBody(domain string, findings []risk.Exposure) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
		<html>
		<head><style>
			body { font-family: Arial, sans-serif; }
			.header { background-color: #dc2626; color: white; padding: 20px; }
			.content { padding: 20px; }
			.finding { border-left: 4px solid #dc2626; padding: 10px; margin: 10px 0; background-color: #fef2f2; }
			.severity { font-weight: bold; color: #dc2626; }
		</style></head>
		<body>
			<div class="header">
				<h1>🚨 Security Alert</h1>
				<p>Domain: <strong>%s</strong></p>
				<p>Critical/High Findings: <strong>%d</strong></p>
			</div>
			<div class="content">
	`, domain, len(findings)))

	for _, f := range findings {
		sb.WriteString(fmt.Sprintf(`
			<div class="finding">
				<h3>%s</h3>
				<p><span class="severity">Severity: %s</span></p>
				<p><strong>Description:</strong> %s</p>
				<p><strong>Remediation:</strong> %s</p>
			</div>
		`, f.Type, f.Severity, f.Description, f.Remediation))
	}

	sb.WriteString(`
			</div>
		</body>
		</html>
	`)

	return sb.String()
}

// FormatEmailBody creates a simple text body for the finding (legacy)
func FormatEmailBody(finding risk.Exposure) string {
	return fmt.Sprintf("Type: %s\nSeverity: %s\nDescription: %s\nRemediation: %s",
		finding.Type, finding.Severity, finding.Description, finding.Remediation)
}
