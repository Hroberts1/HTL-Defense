// Language: Go
package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Note: In order to test with a custom log directory, we temporarily override logDir.
// Since logDir was originally a constant, ensure it is now a variable in production code.
var testLogDir string

func setupTestLogDir(t *testing.T) {
	var err error
	testLogDir, err = ioutil.TempDir("", "testLogs")
	if err != nil {
		t.Fatalf("Failed to create temp log directory: %v", err)
	}
}

func teardownTestLogDir() {
	os.RemoveAll(testLogDir)
}

// TestPruneOldLogs creates one dummy log file with an old modification time and one recent file.
func TestPruneOldLogs(t *testing.T) {
	setupTestLogDir(t)
	defer teardownTestLogDir()

	// Create an "old" log file.
	oldLog := filepath.Join(testLogDir, "old.log")
	if err := ioutil.WriteFile(oldLog, []byte("old log"), 0666); err != nil {
		t.Fatalf("Failed to create old log file: %v", err)
	}
	// Set its mod time to older than logRetentionDays.
	oldTime := time.Now().AddDate(0, 0, -logRetentionDays-1)
	if err := os.Chtimes(oldLog, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set times on old log file: %v", err)
	}

	// Create a "recent" log file.
	recentLog := filepath.Join(testLogDir, "recent.log")
	if err := ioutil.WriteFile(recentLog, []byte("recent log"), 0666); err != nil {
		t.Fatalf("Failed to create recent log file: %v", err)
	}

	// Temporarily override logDir (requires refactoring logDir as a variable in production code).
	origLogDir := logDir
	logDir = testLogDir
	defer func() { logDir = origLogDir }()

	// Run pruneOldLogs
	pruneOldLogs()

	if _, err := os.Stat(oldLog); !os.IsNotExist(err) {
		t.Errorf("Expected old log to be pruned but it still exists")
	}

	if _, err := os.Stat(recentLog); err != nil {
		t.Errorf("Expected recent log to exist but got error: %v", err)
	}
}

// TestRotateLogs simulates a log file exceeding maxLogSize to trigger rotation.
func TestRotateLogs(t *testing.T) {
	// Create a temporary log file.
	tempFile, err := ioutil.TempFile("", "rotate_test_log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create a LogManager for the temp file.
	lm, err := NewLogManager(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create LogManager: %v", err)
	}

	// Override the defaultLogManager with our test instance.
	origDefault := defaultLogManager
	defaultLogManager = lm
	defer func() { defaultLogManager = origDefault }()

	// Write dummy data to exceed maxLogSize.
	dummyData := make([]byte, maxLogSize+100)
	if err := ioutil.WriteFile(tempFile.Name(), dummyData, 0666); err != nil {
		t.Fatalf("Failed to write dummy data: %v", err)
	}

	// Run rotateLogs
	rotateLogs()

	// Stat the current log file; it should be smaller than maxLogSize because a new file was created.
	fi, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat rotated log file: %v", err)
	}

	if fi.Size() >= int64(maxLogSize) {
		t.Errorf("Log file was not rotated; size %d exceeds threshold", fi.Size())
	}
}
