package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"

	"gorm.io/gorm"
)

// SessionPhoneHandler handles GET /sessionphone
func SessionPhoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	sessions := []models.SessionPhone{}
	tx := db.DB
	if sid := q.Get("session_id"); sid != "" {
		tx = tx.Where("session_id = ?", sid)
	}

	// Optional pagination: from/to (1-based, inclusive)
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

	tx.Order("session_id").Find(&sessions)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sessions)
}


// ToggleAIHandler handles PATCH /sessionphone/active
func ToggleAIHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPatch {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    sid := r.URL.Query().Get("session_id")
    result := db.DB.Model(&models.SessionPhone{}).
        Where("session_id = ?", sid).
        Update("ai_active", gorm.Expr("NOT ai_active"))
    if result.RowsAffected == 0 {
        http.Error(w, "session not found", http.StatusNotFound)
        return
    }
    var s models.SessionPhone
    db.DB.First(&s, "session_id = ?", sid)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(s)
}

// SessionDeltaHandler handles GET /sessiondelta
// Returns only sessions with new messages since the given timestamp.
func SessionDeltaHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    sinceParam := r.URL.Query().Get("since")
    var since time.Time
    if sinceParam != "" {
        var err error
        since, err = time.Parse(time.RFC3339, sinceParam)
        if err != nil {
            http.Error(w, "invalid 'since' format, must be RFC3339", http.StatusBadRequest)
            return
        }
    } else {
        // default epoch
        since = time.Unix(0, 0)
    }

    var sessions []models.SessionPhone
    db.DB.
        Where("last_message_at > ?", since).
        Order("last_message_at asc").
        Find(&sessions)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(sessions)
}
