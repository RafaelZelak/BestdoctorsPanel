package routes

import (
    "encoding/json"
    "net/http"

    "bestdoctors_service/internal/db"
    "bestdoctors_service/models"
)

type SessionMetrics struct {
    SessionID                       string  `json:"session_id"`
    DurationSeconds                 float64 `json:"duration_seconds"`
    MessagesCount                   int     `json:"messages_count"`
    AverageInterMessageDelaySeconds float64 `json:"average_inter_message_delay_seconds"`
}

type GlobalMetrics struct {
    TotalSessions                     int64   `json:"total_sessions"`
    AverageDurationSeconds            float64 `json:"average_duration_seconds"`
    AverageMessagesCount              float64 `json:"average_messages_count"`
    AverageInterMessageDelaySeconds   float64 `json:"average_inter_message_delay_seconds"`
}

func SessionMetricsHandler(w http.ResponseWriter, r *http.Request) {
    sessionID := r.URL.Query().Get("session_id")

    if sessionID == "" {
        var ids []string
        if err := db.DB.
            Table(models.SessionPhone{}.TableName()).
            Pluck("session_id", &ids).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        total := int64(len(ids))
        var sumDur, sumMsgs, sumAvgDelays float64

        for _, sid := range ids {
            m, err := computeMetricsForSession(sid)
            if err != nil {
                continue
            }
            sumDur += m.DurationSeconds
            sumMsgs += float64(m.MessagesCount)
            sumAvgDelays += m.AverageInterMessageDelaySeconds
        }

        global := GlobalMetrics{TotalSessions: total}
        if total > 0 {
            global.AverageDurationSeconds = sumDur / float64(total)
            global.AverageMessagesCount = sumMsgs / float64(total)
            global.AverageInterMessageDelaySeconds = sumAvgDelays / float64(total)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(global)
        return
    }

    m, err := computeMetricsForSession(sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}

func computeMetricsForSession(sid string) (SessionMetrics, error) {
    var history []models.ChatHistory
    if err := db.DB.
        Where("session_id = ?", sid).
        Order("created_at ASC").
        Find(&history).Error; err != nil {
        return SessionMetrics{}, err
    }

    count := len(history)
    m := SessionMetrics{SessionID: sid, MessagesCount: count}

    if count > 0 {
        start := history[0].CreatedAt
        end := history[count-1].CreatedAt
        m.DurationSeconds = end.Sub(start).Seconds()
    }
    if count > 1 {
        var sumDelays float64
        for i := 1; i < count; i++ {
            delta := history[i].CreatedAt.Sub(history[i-1].CreatedAt).Seconds()
            if delta > 0 {
                sumDelays += delta
            }
        }
        m.AverageInterMessageDelaySeconds = sumDelays / float64(count-1)
    }
    return m, nil
}
