package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"

	"gorm.io/gorm"
)

type AbandonmentResponse struct {
	TotalSessions          int64   `json:"total_sessions"`
	CompletedSessions      int64   `json:"completed_sessions"`
	AbandonmentRate        float64 `json:"abandonment_rate"`
	TotalEngagedSessions   int64   `json:"total_engaged_sessions"`
	EngagedAbandonmentRate float64 `json:"engaged_abandonment_rate"`
}

type rawMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type contentPayload struct {
	Output struct {
		Vars struct {
			Finalizar bool `json:"finalizar"`
		} `json:"vars"`
	} `json:"output"`
}

// CalculateAbandonmentMetrics calcula a taxa de abandono e métricas relacionadas.
func CalculateAbandonmentMetrics(gormDB *gorm.DB, supabaseDB *gorm.DB) (AbandonmentResponse, error) {
	// Desativa o cache de declaração do GORM para o Supabase para evitar problemas.
	dbNoPrep := supabaseDB.Session(&gorm.Session{PrepareStmt: false})

	var sessionIDs []string
	// A função db.RetryForever não retorna um erro, então simplesmente a chamamos.
	db.RetryForever(100*time.Millisecond, func() error {
		return dbNoPrep.
			Table(models.SessionPhone{}.TableName()).
			Pluck("session_id", &sessionIDs).Error
	})

	totalSessions := int64(len(sessionIDs))
	var completedSessions int64
	var totalEngagedSessions int64

	for _, sid := range sessionIDs {
		var history []models.ChatHistory
		// A função db.RetryForever não retorna um erro, então simplesmente a chamamos.
		db.RetryForever(50*time.Millisecond, func() error {
			return dbNoPrep.
				Where("session_id = ?", sid).
				Order("created_at ASC").
				Find(&history).Error
		})

		if len(history) < 2 {
			continue
		}

		humanCount := 0
		finalizeFlag := false

		for i, entry := range history {
			var rm rawMessage
			if err := json.Unmarshal([]byte(entry.Message), &rm); err != nil {
				continue
			}
			if rm.Type == "human" {
				humanCount++
			}
			if i == len(history)-1 {
				var cp contentPayload
				if err := json.Unmarshal([]byte(rm.Content), &cp); err == nil {
					finalizeFlag = cp.Output.Vars.Finalizar
				}
			}
		}

		if finalizeFlag {
			completedSessions++
		}
		if humanCount > 1 {
			totalEngagedSessions++
		}
	}

	abandonmentRate := 0.0
	if totalSessions > 0 {
		abandonmentRate = float64(totalSessions-completedSessions) / float64(totalSessions) * 100
	}
	engagedAbandonmentRate := 0.0
	if totalEngagedSessions > 0 {
		engagedAbandonmentRate = float64(totalEngagedSessions-completedSessions) / float64(totalEngagedSessions) * 100
	}

	resp := AbandonmentResponse{
		TotalSessions:          totalSessions,
		CompletedSessions:      completedSessions,
		AbandonmentRate:        abandonmentRate,
		TotalEngagedSessions:   totalEngagedSessions,
		EngagedAbandonmentRate: engagedAbandonmentRate,
	}

	return resp, nil
}

func AbandonmentRateHandler(w http.ResponseWriter, r *http.Request) {
	const key = "abandonment"

	// Verificação do cache (24h)
	var cache models.MetricsCache
	if err := db.PostgresDB.
		Where("metric_key = ?", key).
		First(&cache).Error; err == nil {
		if time.Since(cache.LastRefreshedAt) <= 24*time.Hour {
			var resp AbandonmentResponse
			if err := json.Unmarshal(cache.Payload, &resp); err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	resp, err := CalculateAbandonmentMetrics(db.PostgresDB, db.SupabaseDB)
	if err != nil {
		http.Error(w, "Falha ao calcular métricas de abandono", http.StatusInternalServerError)
		return
	}

	// Atualiza o cache
	dataBytes, _ := json.Marshal(resp)
	cache = models.MetricsCache{
		MetricKey:       key,
		Payload:         dataBytes,
		LastRefreshedAt: time.Now(),
	}
	db.PostgresDB.Save(&cache)

	// Retorna a resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
