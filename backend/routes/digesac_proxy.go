package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DigesacProxyHandler forwards all /digesac/homol/* requests to the upstream
// DIGESAC_API_BASE server, injecting Basic Auth credentials from environment.
// This avoids exposing the upstream URL and credentials to the browser.
func DigesacProxyHandler(w http.ResponseWriter, r *http.Request) {
	upstreamBase := os.Getenv("DIGESAC_API_BASE")
	if upstreamBase == "" {
		http.Error(w, "DIGESAC_API_BASE not configured", http.StatusInternalServerError)
		return
	}

	// Strip the frontend prefix to get the upstream path.
	// Incoming:  /digesac/homol/prompts?content=false
	// Upstream:  http://189.45.140.206/api/digesac/homol/prompts?content=false
	upstreamPath := r.URL.Path
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

	proxyReq.SetBasicAuth(
		os.Getenv("DIGESAC_API_USER"),
		os.Getenv("DIGESAC_API_PASS"),
	)

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "upstream request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
