package logging

import (
	"database/sql"
	"sync"
	"time"
)

type SecurityLog struct {
	LogID        string    `json:"log_id"`
	Type         string    `json:"type"`
	Hostname     string    `json:"hostname"`
	HostIP       string    `json:"host_ip"`
	Timestamp    time.Time `json:"timestamp"`
	CVEReference struct {
		CVEID       string `json:"cve_id"`
		Description string `json:"description"`
	} `json:"cve_reference"`
	IncidentReport string `json:"incident_report"`
	ThreatDetails  struct {
		Name            string     `json:"name"`
		Category        string     `json:"category"`
		DetectedAt      time.Time  `json:"detected_at"`
		RemediedAt      *time.Time `json:"remedied_at,omitempty"`
		IncidentCreated time.Time  `json:"incident_created_at"`
		Severity        string     `json:"severity"`
		Action          string     `json:"threat_action"`
		Status          string     `json:"threat_status"`
		DetectionSource string     `json:"detection_source"`
		OSResources     struct {
			PID         int    `json:"pid"`
			ProcessName string `json:"process_name"`
		} `json:"os_resources"`
	} `json:"threat_details"`
	DangerScore int `json:"danger_score"`
}

type Logger struct {
	db    *sql.DB
	mutex sync.Mutex
}

func NewLogger(db *sql.DB) *Logger {
	return &Logger{
		db:    db,
		mutex: sync.Mutex{},
	}
}

func (l *Logger) LogSecurityEvent(event *SecurityLog) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	query := `
        INSERT INTO security_logs (
            type, hostname, host_ip, timestamp, cve_id, cve_description,
            incident_report, threat_name, threat_category, detected_at,
            remedied_at, severity, threat_action, threat_status,
            detection_source, pid, process_name, danger_score
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
        RETURNING log_id`

	return l.db.QueryRow(
		query,
		event.Type,
		event.Hostname,
		event.HostIP,
		event.Timestamp,
		event.CVEReference.CVEID,
		event.CVEReference.Description,
		event.IncidentReport,
		event.ThreatDetails.Name,
		event.ThreatDetails.Category,
		event.ThreatDetails.DetectedAt,
		event.ThreatDetails.RemediedAt,
		event.ThreatDetails.Severity,
		event.ThreatDetails.Action,
		event.ThreatDetails.Status,
		event.ThreatDetails.DetectionSource,
		event.ThreatDetails.OSResources.PID,
		event.ThreatDetails.OSResources.ProcessName,
		event.DangerScore,
	).Scan(&event.LogID)
}

func (l *Logger) LogSecurityEvents(events []*SecurityLog) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO security_logs (
			type, hostname, host_ip, timestamp, cve_id, cve_description,
			incident_report, threat_name, threat_category, detected_at,
			remedied_at, severity, threat_action, threat_status,
			detection_source, pid, process_name, danger_score
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING log_id`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		err := stmt.QueryRow(
			event.Type,
			event.Hostname,
			event.HostIP,
			event.Timestamp,
			event.CVEReference.CVEID,
			event.CVEReference.Description,
			event.IncidentReport,
			event.ThreatDetails.Name,
			event.ThreatDetails.Category,
			event.ThreatDetails.DetectedAt,
			event.ThreatDetails.RemediedAt,
			event.ThreatDetails.Severity,
			event.ThreatDetails.Action,
			event.ThreatDetails.Status,
			event.ThreatDetails.DetectionSource,
			event.ThreatDetails.OSResources.PID,
			event.ThreatDetails.OSResources.ProcessName,
			event.DangerScore,
		).Scan(&event.LogID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
