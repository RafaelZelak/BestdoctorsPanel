package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bestdoctors_service/internal/db"
	"bestdoctors_service/models"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func parseTimeFilter(filters map[string]interface{}) (from *time.Time, to *time.Time) {
	if filters == nil {
		return nil, nil
	}
	if v, ok := filters["from"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = &t
		}
	}
	if v, ok := filters["to"].(string); ok && v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = &t
		}
	}
	return
}

func getSessions(limit ...int) []models.SessionPhone {
	var sessions []models.SessionPhone
	tx := db.DB
	if len(limit) > 0 && limit[0] > 0 {
		tx = tx.Limit(limit[0])
	}
	tx.Order("session_id").Find(&sessions)
	return sessions
}

func getSessionDelta(since time.Time) []models.SessionPhone {
	var sessions []models.SessionPhone
	db.DB.
		Where("last_message_at > ?", since).
		Order("last_message_at asc").
		Find(&sessions)
	return sessions
}

// Sessions com faixa [from, to] aplicada em last_message_at
func getSessionsWithRange(from, to *time.Time) []models.SessionPhone {
	var sessions []models.SessionPhone
	tx := db.DB
	if from != nil {
		tx = tx.Where("last_message_at >= ?", *from)
	}
	if to != nil {
		tx = tx.Where("last_message_at <= ?", *to)
	}
	tx.Order("session_id").Find(&sessions)
	return sessions
}

//
// ──────────────── Cálculo de métricas (com from/to) ────────────────
//
// IMPORTANTE: tipos/funções usadas aqui existem em outros arquivos do mesmo pacote:
// - AbandonmentResponse, rawMessage, contentPayload           → abandonment.go
// - FlowDepthResponse, flowRawMessage, flowContentPayload,
//   detectFlowState                                           → flowdepth.go
// - ReengagementResponse, recapRawMessage                     → reengagement.go
//

// Abandono com filtro por faixa de datas (em last_message_at)
func CalculateAbandonmentMetricsFiltered(supabaseDB *gorm.DB, from, to *time.Time) (AbandonmentResponse, error) {
	dbNoPrep := supabaseDB.Session(&gorm.Session{PrepareStmt: false})

	var sessionIDs []string
	q := dbNoPrep.Table(models.SessionPhone{}.TableName())
	if from != nil {
		q = q.Where("last_message_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("last_message_at <= ?", *to)
	}
	db.RetryForever(100*time.Millisecond, func() error {
		return q.Pluck("session_id", &sessionIDs).Error
	})

	totalSessions := int64(len(sessionIDs))
	var completedSessions int64
	var totalEngagedSessions int64

	for _, sid := range sessionIDs {
		var history []models.ChatHistory
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
		abandonmentRate = float64(totalSessions-completedSessions) / float64(totalSessions) * 100.0
	}
	engagedAbandonmentRate := 0.0
	if totalEngagedSessions > 0 {
		engagedAbandonmentRate = float64(totalEngagedSessions-completedSessions) / float64(totalEngagedSessions) * 100.0
	}

	return AbandonmentResponse{
		TotalSessions:          totalSessions,
		CompletedSessions:      completedSessions,
		AbandonmentRate:        abandonmentRate,
		TotalEngagedSessions:   totalEngagedSessions,
		EngagedAbandonmentRate: engagedAbandonmentRate,
	}, nil
}

// Profundidade do fluxo com filtro por faixa [from, to] (em last_message_at)
func CalculateFlowDepthMetricsFiltered(supabaseDB *gorm.DB, from, to *time.Time) (FlowDepthResponse, error) {
	dbNoPrep := supabaseDB.Session(&gorm.Session{PrepareStmt: false})

	var sessionIDs []string
	q := dbNoPrep.Table(models.SessionPhone{}.TableName())
	if from != nil {
		q = q.Where("last_message_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("last_message_at <= ?", *to)
	}
	db.RetryForever(100*time.Millisecond, func() error {
		return q.Pluck("session_id", &sessionIDs).Error
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
		if total > 0 {
			distributionPercent[state] = float64(count) / float64(total) * 100.0
		} else {
			distributionPercent[state] = 0
		}
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

	return FlowDepthResponse{
		DistributionCount:   depthCount,
		DistributionPercent: distributionPercent,
		AverageDepth:        avgDepth,
		StateLabels:         stateLabels,
	}, nil
}

// Reengajamento com filtro por faixa [from, to] (em last_message_at)
func CalculateReengagementMetricsFiltered(supabaseDB *gorm.DB, includeSessions bool, from, to *time.Time) (ReengagementResponse, error) {
	dbNoPrep := supabaseDB.Session(&gorm.Session{PrepareStmt: false})

	var sessionIDs []string
	q := dbNoPrep.Table(models.SessionPhone{}.TableName())
	if from != nil {
		q = q.Where("last_message_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("last_message_at <= ?", *to)
	}
	db.RetryForever(100*time.Millisecond, func() error {
		return q.Pluck("session_id", &sessionIDs).Error
	})

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
		rate = float64(totalReengaged) / float64(totalRecapture) * 100.0
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
	return resp, nil
}

//
// ─────────────────────────── Report Handler ───────────────────────────
//

type ReportRequest struct {
	Report  string                 `json:"report"` // "session" | "abandonment" | "flowDepth" | "reengagement"
	Type    string                 `json:"type"`   // "json" | "csv" | "pdf" | "xls" | "xlsx"
	Filters map[string]interface{} `json:"filters"`
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	from, to := parseTimeFilter(req.Filters)

	var (
		result interface{}
		err    error
	)

	switch req.Report {
	case "session":
		result = getSessionsWithRange(from, to)

	case "abandonment":
		result, err = CalculateAbandonmentMetricsFiltered(db.SupabaseDB, from, to)

	case "flowDepth":
		result, err = CalculateFlowDepthMetricsFiltered(db.SupabaseDB, from, to)

	case "reengagement":
		include := false
		if req.Filters != nil {
			if v, ok := req.Filters["full"].(bool); ok {
				include = v
			}
		}
		result, err = CalculateReengagementMetricsFiltered(db.SupabaseDB, include, from, to)

	default:
		http.Error(w, "invalid report type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "error generating report: "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Type {
	case "json", "":
		writeAsJSON(w, result)
	case "csv":
		writeAsCSV(w, result, "report.csv")
	case "pdf":
		writeAsPDF(w, result, "report.pdf")
	case "xls", "xlsx":
		writeAsXLSX(w, result, "report.xlsx")
	default:
		http.Error(w, "unsupported export type", http.StatusBadRequest)
	}
}

//
// ──────────────────────────── Saídas ────────────────────────────
//

// JSON
func writeAsJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// CSV (colunas específicas por tipo)
func writeAsCSV(w http.ResponseWriter, data interface{}, filename string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/csv")

	cw := csv.NewWriter(w)
	defer cw.Flush()

	switch v := data.(type) {
	case AbandonmentResponse:
		_ = cw.Write([]string{"total_sessions", "completed_sessions", "abandonment_rate", "total_engaged_sessions", "engaged_abandonment_rate"})
		_ = cw.Write([]string{
			fmt.Sprintf("%d", v.TotalSessions),
			fmt.Sprintf("%d", v.CompletedSessions),
			fmt.Sprintf("%.2f", v.AbandonmentRate),
			fmt.Sprintf("%d", v.TotalEngagedSessions),
			fmt.Sprintf("%.2f", v.EngagedAbandonmentRate),
		})

	case FlowDepthResponse:
		_ = cw.Write([]string{"state", "label", "count", "percent"})
		for state, count := range v.DistributionCount {
			_ = cw.Write([]string{
				fmt.Sprintf("%d", state),
				v.StateLabels[state],
				fmt.Sprintf("%d", count),
				fmt.Sprintf("%.2f", v.DistributionPercent[state]),
			})
		}
		_ = cw.Write([]string{})
		_ = cw.Write([]string{"average_depth", fmt.Sprintf("%.2f", v.AverageDepth)})

	case ReengagementResponse:
		_ = cw.Write([]string{"total_recapture_sessions", "reengaged_sessions", "reengagement_rate"})
		_ = cw.Write([]string{
			fmt.Sprintf("%d", v.TotalRecaptureSessions),
			fmt.Sprintf("%d", v.ReengagedSessions),
			fmt.Sprintf("%.2f", v.ReengagementRate),
		})
		if len(v.RecaptureSessionIDs) > 0 {
			_ = cw.Write([]string{})
			_ = cw.Write([]string{"recapture_session_ids"})
			for _, id := range v.RecaptureSessionIDs {
				_ = cw.Write([]string{id})
			}
		}
		if len(v.ReengagedSessionIDs) > 0 {
			_ = cw.Write([]string{})
			_ = cw.Write([]string{"reengaged_session_ids"})
			for _, id := range v.ReengagedSessionIDs {
				_ = cw.Write([]string{id})
			}
		}

	case []models.SessionPhone:
		_ = cw.Write([]string{
			"session_id", "phone", "ai_active", "created_at", "last_message_at", "lead_name",
		})
		for _, s := range v {
			createdAt := s.CreatedAt.UTC().Format(time.RFC3339Nano)
			lastMsgAt := s.LastMessageAt.UTC().Format(time.RFC3339Nano)
			lead := strings.TrimSpace(s.LeadName)
			if lead == "" {
				lead = "null"
			}
			_ = cw.Write([]string{
				s.SessionID,
				s.Phone,
				fmt.Sprintf("%t", s.AIActive),
				createdAt,
				lastMsgAt,
				lead,
			})
		}

	default:
		// fallback simples
		b, _ := json.Marshal(data)
		_ = cw.Write([]string{"json"})
		_ = cw.Write([]string{string(b)})
	}
}

// XLSX com colunas/abas por tipo
func writeAsXLSX(w http.ResponseWriter, data interface{}, filename string) {
	f := excelize.NewFile()

	writeSheet := func(name string) string {
		if name == "" {
			return f.GetSheetName(f.GetActiveSheetIndex())
		}
		idx, _ := f.NewSheet(name)
		f.SetActiveSheet(idx)
		return name
	}

	switch v := data.(type) {
	case AbandonmentResponse:
		sheet := writeSheet("Abandonment")
		_ = f.SetSheetRow(sheet, "A1", &[]interface{}{"total_sessions", "completed_sessions", "abandonment_rate", "total_engaged_sessions", "engaged_abandonment_rate"})
		_ = f.SetSheetRow(sheet, "A2", &[]interface{}{v.TotalSessions, v.CompletedSessions, fmt.Sprintf("%.2f", v.AbandonmentRate), v.TotalEngagedSessions, fmt.Sprintf("%.2f", v.EngagedAbandonmentRate)})
		_ = f.SetColWidth(sheet, "A", "E", 28)
		_ = f.AutoFilter(sheet, "A1:E2", nil)

	case FlowDepthResponse:
		sheet := writeSheet("FlowDepth")
		_ = f.SetSheetRow(sheet, "A1", &[]interface{}{"state", "label", "count", "percent"})
		row := 2
		for state, count := range v.DistributionCount {
			_ = f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &[]interface{}{state, v.StateLabels[state], count, fmt.Sprintf("%.2f", v.DistributionPercent[state])})
			row++
		}
		_ = f.SetSheetRow(sheet, fmt.Sprintf("A%d", row+1), &[]interface{}{"average_depth", fmt.Sprintf("%.2f", v.AverageDepth)})
		_ = f.SetColWidth(sheet, "A", "D", 26)
		_ = f.AutoFilter(sheet, fmt.Sprintf("A1:D%d", row-1), nil)

	case ReengagementResponse:
		summary := writeSheet("Reengagement")
		_ = f.SetSheetRow(summary, "A1", &[]interface{}{"total_recapture_sessions", "reengaged_sessions", "reengagement_rate"})
		_ = f.SetSheetRow(summary, "A2", &[]interface{}{v.TotalRecaptureSessions, v.ReengagedSessions, fmt.Sprintf("%.2f", v.ReengagementRate)})
		_ = f.SetColWidth(summary, "A", "C", 28)

		if len(v.RecaptureSessionIDs) > 0 {
			s := writeSheet("RecaptureSessions")
			_ = f.SetSheetRow(s, "A1", &[]interface{}{"session_id"})
			for i, id := range v.RecaptureSessionIDs {
				_ = f.SetSheetRow(s, fmt.Sprintf("A%d", i+2), &[]interface{}{id})
			}
			_ = f.SetColWidth(s, "A", "A", 44)
			_ = f.AutoFilter(s, fmt.Sprintf("A1:A%d", len(v.RecaptureSessionIDs)+1), nil)
		}
		if len(v.ReengagedSessionIDs) > 0 {
			s := writeSheet("ReengagedSessions")
			_ = f.SetSheetRow(s, "A1", &[]interface{}{"session_id"})
			for i, id := range v.ReengagedSessionIDs {
				_ = f.SetSheetRow(s, fmt.Sprintf("A%d", i+2), &[]interface{}{id})
			}
			_ = f.SetColWidth(s, "A", "A", 44)
			_ = f.AutoFilter(s, fmt.Sprintf("A1:A%d", len(v.ReengagedSessionIDs)+1), nil)
		}

	case []models.SessionPhone:
		sheet := writeSheet("Sessions")
		_ = f.SetSheetRow(sheet, "A1", &[]interface{}{"session_id", "phone", "ai_active", "created_at", "last_message_at", "lead_name"})
		for i, s := range v {
			row := i + 2
			createdAt := s.CreatedAt.UTC().Format(time.RFC3339Nano)
			lastMsgAt := s.LastMessageAt.UTC().Format(time.RFC3339Nano)
			lead := strings.TrimSpace(s.LeadName)
			if lead == "" {
				lead = "null"
			}
			_ = f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &[]interface{}{
				s.SessionID, s.Phone, s.AIActive, createdAt, lastMsgAt, lead,
			})
		}
		_ = f.SetColWidth(sheet, "A", "F", 26)
		_ = f.AutoFilter(sheet, fmt.Sprintf("A1:F%d", len(v)+1), nil)

	default:
		sheet := writeSheet("Data")
		b, _ := json.MarshalIndent(data, "", "  ")
		_ = f.SetCellValue(sheet, "A1", string(b))
		style, _ := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"}})
		_ = f.SetCellStyle(sheet, "A1", "A1", style)
		_ = f.SetColWidth(sheet, "A", "A", 120)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		http.Error(w, "failed to buffer xlsx", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	_, _ = w.Write(buf.Bytes())
}

// PDF com layout caprichado (tabela para Sessions)
func writeAsPDF(w http.ResponseWriter, data interface{}, filename string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 15, 12)
	pdf.SetAutoPageBreak(true, 12)
	pdf.SetTitle("BestDoctors Report", false)

	addHeader := func(title string) {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 18)
		pdf.Cell(0, 10, "BestDoctors")
		pdf.Ln(9)
		pdf.SetFont("Helvetica", "", 11)
		pdf.SetTextColor(100, 100, 100)
		pdf.Cell(0, 6, time.Now().Format(time.RFC1123))
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(8)
		pdf.SetFont("Helvetica", "B", 16)
		pdf.Cell(0, 9, title)
		pdf.Ln(10)
	}
	writeKV := func(key, val string) {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(60, 7, key, "", 0, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 12)
		pdf.CellFormat(0, 7, val, "", 1, "L", false, 0, "")
	}
	writeKeyValBlock := func(pairs [][2]string) {
		for _, kv := range pairs {
			writeKV(kv[0], kv[1])
		}
	}
	writeSimpleTable := func(headers []string, rows [][]string, colWidths []float64) {
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetFillColor(240, 240, 240)
		for i, h := range headers {
			align := "L"
			if i == len(headers)-1 {
				align = "R"
			}
			pdf.CellFormat(colWidths[i], 7, h, "1", 0, align, true, 0, "")
		}
		pdf.Ln(-1)

		pdf.SetFont("Helvetica", "", 10)
		fill := false
		for _, r := range rows {
			if fill {
				pdf.SetFillColor(248, 248, 248)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			for i, c := range r {
				align := "L"
				if i == len(r)-1 {
					align = "R"
				}
				pdf.CellFormat(colWidths[i], 6.2, c, "1", 0, align, true, 0, "")
			}
			pdf.Ln(-1)
			fill = !fill
		}
	}

	switch v := data.(type) {
	case AbandonmentResponse:
		addHeader("Report: Abandonment")
		writeKeyValBlock([][2]string{
			{"Total Sessions", fmt.Sprintf("%d", v.TotalSessions)},
			{"Completed Sessions", fmt.Sprintf("%d", v.CompletedSessions)},
			{"Abandonment Rate", fmt.Sprintf("%.2f%%", v.AbandonmentRate)},
			{"Total Engaged Sessions", fmt.Sprintf("%d", v.TotalEngagedSessions)},
			{"Engaged Abandonment Rate", fmt.Sprintf("%.2f%%", v.EngagedAbandonmentRate)},
		})

	case FlowDepthResponse:
		addHeader("Report: Flow Depth")
		writeKeyValBlock([][2]string{
			{"Average Depth", fmt.Sprintf("%.2f", v.AverageDepth)},
		})
		pdf.Ln(4)
		headers := []string{"State", "Label", "Count", "Percent"}
		colWidths := []float64{18, 90, 25, 25}
		var rows [][]string
		for state, count := range v.DistributionCount {
			rows = append(rows, []string{
				fmt.Sprintf("%d", state),
				v.StateLabels[state],
				fmt.Sprintf("%d", count),
				fmt.Sprintf("%.2f%%", v.DistributionPercent[state]),
			})
		}
		writeSimpleTable(headers, rows, colWidths)

	case ReengagementResponse:
		addHeader("Report: Reengagement")
		writeKeyValBlock([][2]string{
			{"Total Recapture Sessions", fmt.Sprintf("%d", v.TotalRecaptureSessions)},
			{"Reengaged Sessions", fmt.Sprintf("%d", v.ReengagedSessions)},
			{"Reengagement Rate", fmt.Sprintf("%.2f%%", v.ReengagementRate)},
		})

		if len(v.RecaptureSessionIDs) > 0 {
			pdf.Ln(6)
			pdf.SetFont("Helvetica", "B", 12)
			pdf.Cell(0, 7, "Recapture Session IDs")
			pdf.Ln(8)
			pdf.SetFont("Helvetica", "", 10)
			colW := float64(210-24) / 2.0
			for i, id := range v.RecaptureSessionIDs {
				pdf.CellFormat(colW, 6, id, "", 0, "L", false, 0, "")
				if i%2 == 1 {
					pdf.Ln(6)
				}
			}
			pdf.Ln(4)
		}
		if len(v.ReengagedSessionIDs) > 0 {
			pdf.Ln(4)
			pdf.SetFont("Helvetica", "B", 12)
			pdf.Cell(0, 7, "Reengaged Session IDs")
			pdf.Ln(8)
			pdf.SetFont("Helvetica", "", 10)
			colW := float64(210-24) / 2.0
			for i, id := range v.ReengagedSessionIDs {
				pdf.CellFormat(colW, 6, id, "", 0, "L", false, 0, "")
				if i%2 == 1 {
					pdf.Ln(6)
				}
			}
		}

	case []models.SessionPhone:
		addHeader("Report: Sessions")

		headers := []string{
			"session_id", "phone", "ai_active", "created_at", "last_message_at", "lead_name",
		}
		colWidths := []float64{60, 30, 20, 38, 38, 24}

		writeTableHeader := func() {
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetFillColor(235, 235, 235)
			for i, h := range headers {
				align := "L"
				if h == "ai_active" {
					align = "C"
				}
				pdf.CellFormat(colWidths[i], 7, h, "1", 0, align, true, 0, "")
			}
			pdf.Ln(-1)
		}

		writeTableHeader()

		pdf.SetFont("Helvetica", "", 9)
		fill := false
		rowHeight := 6.0

		for idx, s := range v {
			// quebra de página
			if pdf.GetY() > 280-rowHeight {
				addHeader("Report: Sessions (cont.)")
				writeTableHeader()
			}
			if fill {
				pdf.SetFillColor(248, 248, 248)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			createdAt := s.CreatedAt.UTC().Format(time.RFC3339)
			lastMsgAt := s.LastMessageAt.UTC().Format(time.RFC3339)
			lead := strings.TrimSpace(s.LeadName)
			if lead == "" {
				lead = "null"
			}

			cells := []string{
				s.SessionID,
				s.Phone,
				fmt.Sprintf("%t", s.AIActive),
				createdAt,
				lastMsgAt,
				lead,
			}

			for i, c := range cells {
				align := "L"
				if headers[i] == "ai_active" {
					align = "C"
				}
				pdf.CellFormat(colWidths[i], rowHeight, c, "1", 0, align, true, 0, "")
			}
			pdf.Ln(-1)
			fill = !fill

			if idx >= 2000 {
				pdf.Ln(4)
				pdf.SetFont("Helvetica", "I", 9)
				pdf.Cell(0, 6, "… (truncado)")
				break
			}
		}

	default:
		addHeader("Report: Data")
		pdf.SetFont("Courier", "", 10)
		b, _ := json.MarshalIndent(data, "", "  ")
		pdf.MultiCell(0, 5, string(b), "", "L", false)
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/pdf")
	if err := pdf.Output(w); err != nil {
		http.Error(w, "failed to write pdf", http.StatusInternalServerError)
	}
}
