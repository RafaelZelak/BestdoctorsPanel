package iza_chat

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

type ExtractorPayload struct {
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var (
	extractorPayloads []ExtractorPayload
	extractorMutex    sync.Mutex
)

// ExtractorWebhookReceiverHandler receives any POST and stores the JSON
func ExtractorWebhookReceiverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	extractorMutex.Lock()
	extractorPayloads = append(extractorPayloads, ExtractorPayload{
		Timestamp: time.Now(),
		Data:      payload,
	})
	// Keep only the last 100 to avoid memory leak
	if len(extractorPayloads) > 100 {
		extractorPayloads = extractorPayloads[1:]
	}
	extractorMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ExtractorWebhookPollHandler returns the buffered payloads and clears them
func ExtractorWebhookPollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	extractorMutex.Lock()
	payloads := extractorPayloads
	extractorPayloads = nil // Clear after reading
	extractorMutex.Unlock()

	if payloads == nil {
		payloads = []ExtractorPayload{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payloads)
}
