package database

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// SecurityLogEntry captures detailed security events.
type SecurityLogEntry struct {
	LogID          string       `json:"log_id"`
	Type           string       `json:"type"`
	Hostname       string       `json:"hostname"`
	HostIP         string       `json:"host_ip"`
	Timestamp      string       `json:"timestamp"`
	CVEReference   string       `json:"cve_reference"`
	IncidentReport string       `json:"incident_report"`
	ThreatDetails  ThreatDetail `json:"threat_details"`
}

type ThreatDetail struct {
	ThreatName              string     `json:"threat_name"`
	ThreatCategory          string     `json:"threat_category"`
	DetectedAt              string     `json:"detected_at"`
	RemediedAt              string     `json:"remedied_at"`
	IncidentReportCreatedAt string     `json:"incident_report_created_at"`
	Severity                string     `json:"severity"`
	ThreatAction            string     `json:"threat_action"`
	ThreatStatus            string     `json:"threat_status"`
	DetectionSource         string     `json:"detection_source"`
	OSResources             OSResource `json:"os_resources"`
}

type OSResource struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"process_name"`
}

type LogManager struct {
	mu       sync.Mutex
	logs     []LogEntry
	logFile  *os.File
	filePath string
}

var logDir = "../logs"

const logRetentionDays = 90
const maxLogSize = 10 * 1024 * 1024 // 10 MB

var defaultLogManager *LogManager

// NewLogManager opens/creates the log file and returns a LogManager.
func NewLogManager(filePath string) (*LogManager, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	lm := &LogManager{
		logFile:  file,
		filePath: filePath,
		logs:     []LogEntry{},
	}
	defaultLogManager = lm
	return lm, nil
}

// AddLog adds a standard log entry.
func (lm *LogManager) AddLog(entry LogEntry) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.logs = append(lm.logs, entry)

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = lm.logFile.Write(data)
	if err != nil {
		return err
	}

	_, err = lm.logFile.WriteString("\n")
	return err
}

// AddSecurityLog adds a security log entry.
func (lm *LogManager) AddSecurityLog(entry SecurityLogEntry) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = lm.logFile.Write(data)
	if err != nil {
		return err
	}
	_, err = lm.logFile.WriteString("\n")
	return err
}

// Close shuts down the LogManager.
func (lm *LogManager) Close() error {
	return lm.logFile.Close()
}

// rotateLogs checks the current log file size and rotates if necessary.
func rotateLogs() {
	if defaultLogManager == nil {
		return
	}

	defaultLogManager.mu.Lock()
	defer defaultLogManager.mu.Unlock()

	fileInfo, err := os.Stat(defaultLogManager.filePath)
	if err != nil {
		log.Printf("Error stating log file: %v", err)
		return
	}

	if fileInfo.Size() < maxLogSize {
		return
	}

	// Close current log file.
	err = defaultLogManager.logFile.Close()
	if err != nil {
		log.Printf("Error closing log file: %v", err)
		return
	}

	// Add a short delay to allow Windows to release the file handle.
	time.Sleep(50 * time.Millisecond)

	// Rename the current log file with a timestamp.
	archiveName := fmt.Sprintf("%s.%s", defaultLogManager.filePath, time.Now().Format("20060102T150405"))
	err = os.Rename(defaultLogManager.filePath, archiveName)
	if err != nil {
		log.Printf("Error renaming log file: %v", err)
		return
	}

	// Open a new log file.
	newFile, err := os.OpenFile(defaultLogManager.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error creating new log file: %v", err)
		return
	}
	defaultLogManager.logFile = newFile
	log.Println("Log file rotated.")
}

// pruneOldLogs deletes log files older than logRetentionDays in the log directory.
func pruneOldLogs() {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		log.Printf("Error reading log directory: %v", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -logRetentionDays)
	for _, file := range files {
		if file.Mode().IsRegular() {
			fullPath := filepath.Join(logDir, file.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Printf("Error stating file %s: %v", fullPath, err)
				continue
			}
			if info.ModTime().Before(cutoff) {
				err = os.Remove(fullPath)
				if err != nil {
					log.Printf("Error removing file %s: %v", fullPath, err)
				} else {
					log.Printf("Pruned old log file: %s", fullPath)
				}
			}
		}
	}
}

func startLogMaintenance() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		log.Println("Starting maintenance: pruning and rotating logs.")
		pruneOldLogs()
		rotateLogs()
	}
}

func init() {
	// Ensure the log directory exists.
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, fs.ModePerm)
		if err != nil {
			log.Fatalf("Error creating log directory: %v", err)
		}
	}
	go startLogMaintenance()
}
