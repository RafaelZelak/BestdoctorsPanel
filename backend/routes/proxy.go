package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// GenericProxyHandler creates a reverse proxy handler that forwards requests to an upstream server.
// It dynamically resolves the upstream URL and credentials using the provided environment variable keys.
// The stripPrefix is removed from the path before appending it to the base upstream URL.
func GenericProxyHandler(envBaseURLKey, envUserKey, envPassKey, stripPrefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upstreamBase := os.Getenv(envBaseURLKey)
		if upstreamBase == "" {
			http.Error(w, fmt.Sprintf("%s not configured", envBaseURLKey), http.StatusInternalServerError)
			return
		}

		// Strip the specified prefix so only the tail appends to the upstream Base URL
		upstreamPath := strings.TrimPrefix(r.URL.Path, stripPrefix)
		if !strings.HasPrefix(upstreamPath, "/") {
			upstreamPath = "/" + upstreamPath
		}
		upstreamURL := fmt.Sprintf("%s%s", strings.TrimRight(upstreamBase, "/"), upstreamPath)
		if r.URL.RawQuery != "" {
			upstreamURL += "?" + r.URL.RawQuery
		}

		proxyReq, err := http.NewRequest(r.Method, upstreamURL, r.Body)
		if err != nil {
			http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
			return
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "" {
			proxyReq.Header.Set("Content-Type", contentType)
		}

		// Use auth if expected/provided in environment
		apiUser := os.Getenv(envUserKey)
		apiPass := os.Getenv(envPassKey)
		if apiUser != "" || apiPass != "" {
			proxyReq.SetBasicAuth(apiUser, apiPass)
		}

		// Use a custom HTTP client with a strict timeout to prevent Goroutine leaks
		// if the upstream server hangs indefinitely.
		client := &http.Client{
			// Using 30 seconds as max timeout for AI prompt generations which might take a bit
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "upstream request failed or timed out", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
