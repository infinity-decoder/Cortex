package middleware

import (
	"net/http"
	"os"
	"strings"
)

// HTTPSRedirect middleware enforces HTTPS in production
func HTTPSRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce HTTPS if ENFORCE_HTTPS is set to "true"
		if os.Getenv("ENFORCE_HTTPS") == "true" {
			// Check if request is already HTTPS
			proto := r.Header.Get("X-Forwarded-Proto")
			if proto != "https" {
				// Redirect to HTTPS
				httpsURL := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// HSTS adds HTTP Strict Transport Security header
func HSTS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only add HSTS in production with HTTPS
		if os.Getenv("ENFORCE_HTTPS") == "true" {
			proto := r.Header.Get("X-Forwarded-Proto")
			if proto == "https" || strings.HasPrefix(r.URL.Scheme, "https") {
				// HSTS: max-age=31536000 (1 year), includeSubDomains, preload
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
		}
		next.ServeHTTP(w, r)
	})
}
