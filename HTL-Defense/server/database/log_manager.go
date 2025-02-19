package database

import (
    "encoding/json"
    "os"
    "sync"
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