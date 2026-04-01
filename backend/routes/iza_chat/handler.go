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
	"time"

	"github.com/google/uuid"
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

var jsonCodeBlockPattern = regexp.MustCompile("(?s)```(?:json)?\\s*([\\s\\S]*?)\\s*```")

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
	httpClient := &http.Client{Timeout: 60 * time.Second}

	upstreamResponse, err := httpClient.Do(upstreamRequest)
	if err != nil {
		http.Error(w, "upstream request failed or timed out", http.StatusBadGateway)
		return
	}
	defer upstreamResponse.Body.Close()

	responseBody, err := io.ReadAll(upstreamResponse.Body)
	if err != nil {
		http.Error(w, "failed to read upstream response", http.StatusInternalServerError)
		return
	}

	var izaResponse izaWebhookResponse
	if err := json.Unmarshal(responseBody, &izaResponse); err != nil {
		http.Error(w, "failed to parse upstream response", http.StatusInternalServerError)
		return
	}

	agentMessage, err := extractAgentMessage(izaResponse.Response)
	if err != nil {
		http.Error(w, "failed to extract agent message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatMessageResponse{Message: agentMessage})
}
