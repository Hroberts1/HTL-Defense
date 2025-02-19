package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type LogManager struct {
	mu      sync.Mutex
	logs    []LogEntry
	logFile *os.File
}

const logDir = "../logs"
const logRetentionDays = 90
const maxLogSize = 10 * 1024 * 1024 // 10 MB

func NewLogManager(filePath string) (*LogManager, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &LogManager{
		logFile: file,
		logs:    []LogEntry{},
	}, nil
}

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

func (lm *LogManager) Close() error {
	return lm.logFile.Close()
}

func init() {
	go startLogMaintenance()
}

func startLogMaintenance() {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		<-ticker.C
		pruneOldLogs()
		rotateLogs()
	}
}

func pruneOldLogs() {
	fmt.Println("Pruning logs older than 90 days...")
	// TODO: Implement deletion of logs older than 90 days.
}

func rotateLogs() {
	fmt.Println("Rotating logs to archive old logs...")
	// TODO: Implement log rotation logic (e.g., archive current log file and start a new one).
}
