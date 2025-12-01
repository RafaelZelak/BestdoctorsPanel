package routes

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "bestdoctors_service/internal/db"
    "bestdoctors_service/models"
    "gorm.io/gorm"
)

type ReengagementResponse struct {
    TotalRecaptureSessions int64    `json:"total_recapture_sessions"`
    RecaptureSessionIDs    []string `json:"recapture_session_ids,omitempty"`
    ReengagedSessions      int64    `json:"reengaged_sessions"`
    ReengagedSessionIDs    []string `json:"reengaged_session_ids,omitempty"`
    ReengagementRate       float64  `json:"reengagement_rate"`
}

type recapRawMessage struct {
    Type    string `json:"type"`
    Content string `json:"content"`
}

func ReengagementRateHandler(w http.ResponseWriter, r *http.Request) {
    const key = "reengagement"

    var cache models.MetricsCache
    if err := db.PostgresDB.
        Where("metric_key = ?", key).
        First(&cache).Error; err == nil {
        if time.Since(cache.LastRefreshedAt) <= 24*time.Hour {
            var resp ReengagementResponse
            if err := json.Unmarshal(cache.Payload, &resp); err == nil {
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(resp)
                return
            }
        }
    }

    dbNoPrep := db.SupabaseDB.Session(&gorm.Session{PrepareStmt: false})

    var sessionIDs []string
    db.RetryForever(100*time.Millisecond, func() error {
        return dbNoPrep.
            Table(models.SessionPhone{}.TableName()).
            Pluck("session_id", &sessionIDs).Error
    })

    includeSessions := r.URL.Query().Get("sessions") == "true"
    recaptureSet := make(map[string]struct{})
    reengagedSet := make(map[string]struct{})

    for _, sid := range sessionIDs {
        var history []models.ChatHistory
        db.RetryForever(50*time.Millisecond, func() error {
            return dbNoPrep.
                Where("session_id = ?", sid).
                Order("created_at ASC").
                Find(&history).Error
        })

        recaptureIdx := -1
        for i, entry := range history {
            var rm recapRawMessage
            if err := json.Unmarshal([]byte(entry.Message), &rm); err != nil {
                continue
            }
            if rm.Type == "human" && strings.HasPrefix(rm.Content, "Recapture - ") {
                recaptureIdx = i
                break
            }
        }
        if recaptureIdx < 0 {
            continue
        }
        recaptureSet[sid] = struct{}{}

        for j := recaptureIdx + 1; j < len(history); j++ {
            var rm recapRawMessage
            if err := json.Unmarshal([]byte(history[j].Message), &rm); err != nil {
                continue
            }
            if rm.Type == "human" && !strings.HasPrefix(rm.Content, "Recapture - ") {
                reengagedSet[sid] = struct{}{}
                break
            }
        }
    }

    recaptureSessionIDs := make([]string, 0, len(recaptureSet))
    for sid := range recaptureSet {
        recaptureSessionIDs = append(recaptureSessionIDs, sid)
    }
    reengagedSessionIDs := make([]string, 0, len(reengagedSet))
    for sid := range reengagedSet {
        reengagedSessionIDs = append(reengagedSessionIDs, sid)
    }

    totalRecapture := int64(len(recaptureSessionIDs))
    totalReengaged := int64(len(reengagedSessionIDs))
    rate := 0.0
    if totalRecapture > 0 {
        rate = float64(totalReengaged) / float64(totalRecapture) * 100
    }

    resp := ReengagementResponse{
        TotalRecaptureSessions: totalRecapture,
        ReengagedSessions:      totalReengaged,
        ReengagementRate:       rate,
    }
    if includeSessions {
        resp.RecaptureSessionIDs = recaptureSessionIDs
        resp.ReengagedSessionIDs = reengagedSessionIDs
    }

    dataBytes, _ := json.Marshal(resp)
    cache = models.MetricsCache{
        MetricKey:       key,
        Payload:         dataBytes,
        LastRefreshedAt: time.Now(),
    }
    db.PostgresDB.Save(&cache)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
