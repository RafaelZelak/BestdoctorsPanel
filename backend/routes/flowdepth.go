// routes/flowdepth.go
package routes

import (
    "encoding/json"
    "net/http"
    "time"

    "bestdoctors_service/internal/db"
    "bestdoctors_service/models"
    "gorm.io/gorm"
)

type FlowDepthResponse struct {
    DistributionCount   map[int]int64   `json:"distribution_count"`
    DistributionPercent map[int]float64 `json:"distribution_percent"`
    AverageDepth        float64         `json:"average_depth"`
    StateLabels         map[int]string  `json:"state_labels"`
}

type flowRawMessage struct {
    Type    string `json:"type"`
    Content string `json:"content"`
}

type flowContentPayload struct {
    Output struct {
        Vars struct {
            SaudacaoEnviada  bool   `json:"saudacao_enviada"`
            Especialidade    string `json:"especialidade"`
            NomeDoLead       string `json:"nome_do_lead"`
            NumeroDeUsuarios int    `json:"numero_de_usuarios"`
            Finalizar        bool   `json:"finalizar"`
        } `json:"vars"`
    } `json:"output"`
}

func detectFlowState(cp flowContentPayload) int {
    if cp.Output.Vars.Finalizar {
        return 5
    }
    if cp.Output.Vars.NumeroDeUsuarios > 0 {
        return 4
    }
    if cp.Output.Vars.NomeDoLead != "" {
        return 3
    }
    if cp.Output.Vars.Especialidade != "" {
        return 2
    }
    if cp.Output.Vars.SaudacaoEnviada {
        return 1
    }
    return 0
}

func FlowDepthHandler(w http.ResponseWriter, r *http.Request) {
    const key = "flowdepth"

    var cache models.MetricsCache
    if err := db.PostgresDB.
        Where("metric_key = ?", key).
        First(&cache).Error; err == nil {
        if time.Since(cache.LastRefreshedAt) <= 24*time.Hour {
            var resp FlowDepthResponse
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

    total := int64(len(sessionIDs))
    depthCount := make(map[int]int64, 6)
    var sumDepth int64

    for _, sid := range sessionIDs {
        var history []models.ChatHistory
        db.RetryForever(50*time.Millisecond, func() error {
            return dbNoPrep.
                Where("session_id = ?", sid).
                Order("created_at ASC").
                Find(&history).Error
        })

        if len(history) == 0 {
            depthCount[0]++
            continue
        }

        maxState := 0
        for _, entry := range history {
            var rm flowRawMessage
            if err := json.Unmarshal([]byte(entry.Message), &rm); err != nil {
                continue
            }
            if rm.Type != "ai" {
                continue
            }
            var cp flowContentPayload
            if err := json.Unmarshal([]byte(rm.Content), &cp); err != nil {
                continue
            }
            st := detectFlowState(cp)
            if st > maxState {
                maxState = st
                if st == 5 {
                    break
                }
            }
        }

        depthCount[maxState]++
        sumDepth += int64(maxState)
    }

    distributionPercent := make(map[int]float64, len(depthCount))
    for state, count := range depthCount {
        distributionPercent[state] = float64(count) / float64(total) * 100
    }
    avgDepth := 0.0
    if total > 0 {
        avgDepth = float64(sumDepth) / float64(total)
    }

    stateLabels := map[int]string{
        0: "no_ai_response",
        1: "saudacao_enviada",
        2: "especialidade_informada",
        3: "nome_do_lead_informado",
        4: "numero_de_usuarios_informado",
        5: "finalizar_true",
    }

    resp := FlowDepthResponse{
        DistributionCount:   depthCount,
        DistributionPercent: distributionPercent,
        AverageDepth:        avgDepth,
        StateLabels:         stateLabels,
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
