package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type LogsHandler struct {
	db *sql.DB
}

type SecurityLog struct {
	LogID         int       `json:"log_id"`
	Type          string    `json:"type"`
	Hostname      string    `json:"hostname"`
	HostIP        string    `json:"host_ip"`
	Timestamp     string    `json:"timestamp"`
	CVEReference  struct {
		CVEID       string `json:"cve_id"`
		Description string `json:"description"`
	} `json:"cve_reference"`
	IncidentReport string `json:"incident_report"`
	ThreatDetails  struct {
		Name           string `json:"name"`
		Category       string `json:"category"`
		DetectedAt     string `json:"detected_at"`
		RemediedAt     string `json:"remedied_at"`
		IncidentCreated string `json:"incident_created"`
		Severity       string `json:"severity"`
		Action         string `json:"action"`
		Status         string `json:"status"`
		DetectionSource string `json:"detection_source"`
		OSResources    struct {
			PID         int    `json:"pid"`
			ProcessName string `json:"process_name"`
		} `json:"os_resources"`
	} `json:"threat_details"`
	DangerScore float64 `json:"danger_score"`
}

func NewLogsHandler(db *sql.DB) *LogsHandler {
	return &LogsHandler{db: db}
}

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	// Implement rate limiting
	return func(w http.ResponseWriter, r *http.Request) {
		// Add rate limiting logic here
		next.ServeHTTP(w, r)
	}
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	// Implement authentication
	return func(w http.ResponseWriter, r *http.Request) {
		// Add authentication logic here
		next.ServeHTTP(w, r)
	}
}

func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	severity := r.URL.Query().Get("severity")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	category := r.URL.Query().Get("category")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Set default pagination values
	pageNum := 1
	limitNum := 50
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}
	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitNum = l
	}

	// Build query with filters
	query := `
        SELECT log_id, type, hostname, host_ip, timestamp, 
            cve_id, cve_description, incident_report,
            threat_name, threat_category, detected_at, remedied_at,
            incident_created_at, severity, threat_action, threat_status,
            detection_source, pid, process_name, danger_score
        FROM security_logs
        WHERE ($1 = '' OR severity = $1)
        AND ($2 = '' OR timestamp >= $2::timestamp)
        AND ($3 = '' OR timestamp <= $3::timestamp)
        AND ($4 = '' OR threat_category = $4)
        ORDER BY timestamp DESC
        LIMIT $5 OFFSET $6`

	offset := (pageNum - 1) * limitNum
	rows, err := h.db.Query(query, severity, startDate, endDate, category, limitNum, offset)
	if err != nil {
		http.Error(w, "Error querying logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []SecurityLog
	for rows.Next() {
		var log SecurityLog
		err := rows.Scan(
			&log.LogID, &log.Type, &log.Hostname, &log.HostIP, &log.Timestamp,
			&log.CVEReference.CVEID, &log.CVEReference.Description, &log.IncidentReport,
			&log.ThreatDetails.Name, &log.ThreatDetails.Category, &log.ThreatDetails.DetectedAt,
			&log.ThreatDetails.RemediedAt, &log.ThreatDetails.IncidentCreated, &log.ThreatDetails.Severity,
			&log.ThreatDetails.Action, &log.ThreatDetails.Status, &log.ThreatDetails.DetectionSource,
			&log.ThreatDetails.OSResources.PID, &log.ThreatDetails.OSResources.ProcessName, &log.DangerScore,
		)
		if err != nil {
			http.Error(w, "Error scanning logs", http.StatusInternalServerError)
			return
		}
		logs = append(logs, log)
	}

	response := map[string]interface{}{
		"logs": logs,
		"pagination": map[string]int{
			"page":  pageNum,
			"limit": limitNum,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
