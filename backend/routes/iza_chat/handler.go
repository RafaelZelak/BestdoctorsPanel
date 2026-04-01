package iza_chat

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	responseChannels = make(map[string]chan string)
	responseMutex    sync.RWMutex
	jsonCodeBlockPattern = regexp.MustCompile(`(?s)` + "```" + `(?:json)?\s*([\s\S]*?)\s*` + "```" + ``)
)

const (
	izaChatOrigin            = "homol"
	izaChatVersion           = "V2"
	izaChatClientID          = 1
	izaChatLawyerID          = 11
	izaChatFirstInteraction  = false
	izaChatConversationTTL   = 3600000
)

type chatMessageRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id"`
}

type izaWebhookPayload struct {
	Message              string `json:"message"`
	ConversationID       string `json:"conversation_id"`
	ClientID             int    `json:"client_id"`
	ClientName           string `json:"client_name"`
	LawyerID             int    `json:"lawyer_id"`
	FirstInteraction     bool   `json:"first_interaction"`
	ConversationLifetime int    `json:"conversation_lifetime"`
}

type izaWebhookResponse struct {
	Response string `json:"response"`
}

type izaAgentMessage struct {
	Message    string `json:"message"`
	FinishType string `json:"finish_type"`
}

type chatMessageResponse struct {
	Message string `json:"message"`
}

func generateSessionClientName(conversationID string) string {
	hash := md5.Sum([]byte(conversationID))
	return fmt.Sprintf("TEST-%x", hash[:6])
}


func extractAgentMessage(rawResponse string) (string, error) {
	trimmed := strings.TrimSpace(rawResponse)

	jsonContent := trimmed
	if matches := jsonCodeBlockPattern.FindStringSubmatch(trimmed); len(matches) >= 2 {
		jsonContent = strings.TrimSpace(matches[1])
	}

	var agentMessage izaAgentMessage
	if err := json.Unmarshal([]byte(jsonContent), &agentMessage); err != nil {
		return trimmed, nil
	}

	return agentMessage.Message, nil
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var chatRequest chatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(chatRequest.Message) == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(chatRequest.ConversationID) == "" {
		chatRequest.ConversationID = uuid.New().String()
	}

	izaChatURL := os.Getenv("IZA_CHAT_URL")
	if izaChatURL == "" {
		http.Error(w, "IZA_CHAT_URL not configured", http.StatusInternalServerError)
		return
	}

	clientName := generateSessionClientName(chatRequest.ConversationID)

	webhookPayload := izaWebhookPayload{
		Message:              chatRequest.Message,
		ConversationID:       chatRequest.ConversationID,
		ClientID:             izaChatClientID,
		ClientName:           clientName,
		LawyerID:             izaChatLawyerID,
		FirstInteraction:     izaChatFirstInteraction,
		ConversationLifetime: izaChatConversationTTL,
	}

	payloadBytes, err := json.Marshal(webhookPayload)
	if err != nil {
		http.Error(w, "failed to build request payload", http.StatusInternalServerError)
		return
	}

	upstreamURL := fmt.Sprintf("%s?origin=%s&version=%s",
		strings.TrimRight(izaChatURL, "/"),
		izaChatOrigin,
		izaChatVersion,
	)

	upstreamRequest, err := http.NewRequest(http.MethodPost, upstreamURL, bytes.NewReader(payloadBytes))
	if err != nil {
		http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
		return
	}
	upstreamRequest.Header.Set("Content-Type", "application/json")

	// IZA AI responses can take several seconds to generate
	// Fire and forget the request (with a small timeout for the initial connection)
	go func() {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(upstreamRequest)
		if err == nil {
			resp.Body.Close()
		}
	}()

	// Create a channel to wait for the webhook response
	ch := make(chan string, 1)
	responseMutex.Lock()
	responseChannels[chatRequest.ConversationID] = ch
	responseMutex.Unlock()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "processing"})
}

// WebhookReceiverHandler is called by the IZA agent when it finishes processing
func WebhookReceiverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	var izaResponse izaWebhookResponse
	if err := json.Unmarshal(bodyBytes, &izaResponse); err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	// Some payloads might send generic webhook data, try to parse it specifically
	// if we don't have response inside 'response' top-level key.
	// We'll extract conversation_id first to route the message.
	var genericPayload map[string]interface{}
	json.Unmarshal(bodyBytes, &genericPayload)
	convIDRaw, ok := genericPayload["conversation_id"]
	if !ok || convIDRaw == "" {
		http.Error(w, "missing conversation_id", http.StatusBadRequest)
		return
	}
	convID := fmt.Sprintf("%v", convIDRaw)
	
	agentMessage, _ := extractAgentMessage(izaResponse.Response)

	// Send message to waiting client
	responseMutex.Lock()
	if ch, exists := responseChannels[convID]; exists {
		ch <- agentMessage
		delete(responseChannels, convID)
	}
	responseMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// PollResponseHandler is called by the frontend to long-poll for the response
func PollResponseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	convID := r.URL.Query().Get("conversation_id")
	if convID == "" {
		http.Error(w, "missing conversation_id", http.StatusBadRequest)
		return
	}

	responseMutex.RLock()
	ch, exists := responseChannels[convID]
	responseMutex.RUnlock()

	if !exists {
		// No pending response / already answered / wrong ID
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"done": false})
		return
	}

	// Wait for the message with a timeout of 120s (long polling)
	select {
	case msg := <-ch:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatMessageResponse{Message: msg})
	case <-time.After(120 * time.Second):
		// Remove from map
		responseMutex.Lock()
		delete(responseChannels, convID)
		responseMutex.Unlock()
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"done": false, "error": "timeout"})
	}
}
