package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"
)

type SendMessageRequest struct {
	To        string                 `json:"to"`
	Message   string                 `json:"message"`
	SessionID string                 `json:"session_id"`
	Vars      map[string]interface{} `json:"vars"`
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Helper para limpar CRLF (Windows line endings)
	cleanEnv := func(key string) string {
		v := os.Getenv(key)
		v = strings.TrimSpace(v)
		v = strings.ReplaceAll(v, "\r", "")
		v = strings.ReplaceAll(v, "\n", "")
		return v
	}

	// --- pega as vars da .env ---
	accountSID := cleanEnv("TWILIO_ACCOUNT_SID")
	authToken := cleanEnv("TWILIO_AUTH_TOKEN")
	apiURL := cleanEnv("TWILIO_URL")
	fromNumber := "whatsapp:+554134111916"

	// Validação básica
	if accountSID == "" || authToken == "" || apiURL == "" {
		log.Println("❌ Twilio credentials missing in .env")
		http.Error(w, "Twilio configuration missing", http.StatusInternalServerError)
		return
	}

	// --- envia mensagem pelo Twilio ---
	form := url.Values{}
	form.Set("To", req.To)
	form.Set("From", fromNumber)
	form.Set("Body", req.Message)

	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("❌ Failed to create Twilio request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	
	httpReq.SetBasicAuth(accountSID, authToken)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("twilio request failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		http.Error(w, fmt.Sprintf("twilio error %d: %s", resp.StatusCode, string(body)), http.StatusBadGateway)
		return
	}

	// --- monta o JSON do content como STRING ---
	output := map[string]interface{}{
		"output": map[string]interface{}{
			"message": req.Message,
			"vars":    req.Vars,
		},
	}
	outputJSON, _ := json.MarshalIndent(output, "", "  ")

	// Objeto final a ser salvo no banco
	msg := map[string]interface{}{
		"type":               "ai",
		"content":            string(outputJSON), // vira string JSON escapada
		"tool_calls":         []interface{}{},
		"additional_kwargs":  map[string]interface{}{},
		"response_metadata":  map[string]interface{}{},
		"invalid_tool_calls": []interface{}{},
	}
	jsonMsg, _ := json.Marshal(msg)

	// --- salva no banco ---
	history := models.ChatHistory{
		SessionID: req.SessionID,
		Message:   string(jsonMsg),
		CreatedAt: time.Now().UTC(),
	}
	db.DB.Create(&history)

	// --- responde ao cliente ---
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sent",
		"twilio":  json.RawMessage(body),
		"history": history,
	})
}
