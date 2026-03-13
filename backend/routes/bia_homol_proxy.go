package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// BiaHomolProxyHandler forwards all /bia/homol/* requests to the upstream
// BIA_HOMOL_URL server, optionally injecting Basic Auth credentials from environment
// if they are provided, similar to the digesac proxy.
func BiaHomolProxyHandler(w http.ResponseWriter, r *http.Request) {
	upstreamBase := os.Getenv("BIA_HOMOL_URL")
	if upstreamBase == "" {
		http.Error(w, "BIA_HOMOL_URL not configured", http.StatusInternalServerError)
		return
	}

	// Strip the /bia/homol prefix so only the tail appends to BIA_HOMOL_URL.
	upstreamPath := strings.TrimPrefix(r.URL.Path, "/bia/homol")
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

	// Use auth if provided in environment
	apiUser := os.Getenv("BIA_HOMOL_API_USER")
	apiPass := os.Getenv("BIA_HOMOL_API_PASS")
	if apiUser != "" || apiPass != "" {
		proxyReq.SetBasicAuth(apiUser, apiPass)
	}

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
