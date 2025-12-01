package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"
)


type HistoryResponse struct {
	ID         uint        `json:"id"`
	SessionID  string      `json:"session_id"`
	CreatedAt  time.Time   `json:"created_at"`
	Message    interface{} `json:"message"`               
	MessageRaw *string     `json:"message_raw,omitempty"` 
}

func parseMessage(raw string) interface{} {
	var top map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &top); err != nil {
		return raw 
	}
	if c, ok := top["content"]; ok {
		if s, isString := c.(string); isString {
			var nested interface{}
			if err := json.Unmarshal([]byte(s), &nested); err == nil {
				top["content"] = nested
			}
		}
	}
	return top
}

func ChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHistoryNormalized(w, r)
	case http.MethodPost:
		postHistory(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getHistoryNormalized(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sid := q.Get("session_id")
	if sid == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}

	tx := db.DB.Where("session_id = ?", sid)
	if since := q.Get("since"); since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			tx = tx.Where("created_at > ?", t)
		}
	}

	var (
		fromVal int
		toVal   int
	)
	if v := q.Get("from"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			fromVal = n
		}
	}
	if v := q.Get("to"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			toVal = n
		}
	}

	if fromVal > 0 || toVal > 0 {
		if fromVal <= 0 {
			fromVal = 1
		}
		offset := fromVal - 1
		if toVal > 0 && toVal >= fromVal {
			limit := toVal - offset
			tx = tx.Offset(offset).Limit(limit)
		} else {
			tx = tx.Offset(offset)
		}
	}

	var rows []models.ChatHistory
	tx.Order("created_at asc").Find(&rows)

	includeRaw := q.Get("raw") == "1" || q.Get("raw") == "true"
	pretty := q.Get("pretty") == "1" || q.Get("pretty") == "true"

	out := make([]HistoryResponse, 0, len(rows))
	for _, h := range rows {
		item := HistoryResponse{
			ID:        h.ID,
			SessionID: h.SessionID,
			CreatedAt: h.CreatedAt,
			Message:   parseMessage(h.Message),
		}
		if includeRaw {
			raw := h.Message
			item.MessageRaw = &raw
		}
		out = append(out, item)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if pretty {
		blob, _ := json.MarshalIndent(out, "", "  ")
		_, _ = w.Write(blob)
		return
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(out)
}

func postHistory(w http.ResponseWriter, r *http.Request) {
	var h models.ChatHistory
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	db.DB.Create(&h)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(h)
}
