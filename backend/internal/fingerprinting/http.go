package fingerprinting

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Fingerprint struct {
	Server      string
	PoweredBy   string
	BodySnippet string
	StatusCode  int
}

// HTTPFingerprint basic HTTP response analysis
func HTTPFingerprint(ctx context.Context, url string) (*Fingerprint, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: nil, // Add insecure if needed for internal dev
		},
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Read first 1KB

	fp := &Fingerprint{
		Server:      resp.Header.Get("Server"),
		PoweredBy:   resp.Header.Get("X-Powered-By"),
		BodySnippet: string(body),
		StatusCode:  resp.StatusCode,
	}

	return fp, nil
}
