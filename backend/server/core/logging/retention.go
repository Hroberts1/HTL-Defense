package logging

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type RetentionManager struct {
	db            *sql.DB
	retentionDays int
}

func NewRetentionManager(db *sql.DB, retentionDays int) *RetentionManager {
	return &RetentionManager{
		db:            db,
		retentionDays: retentionDays,
	}
}

func (rm *RetentionManager) StartRetentionPolicy(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := rm.pruneOldLogs(); err != nil {
					log.Printf("Error pruning logs: %v", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (rm *RetentionManager) pruneOldLogs() error {
	query := `DELETE FROM security_logs WHERE timestamp < NOW() - INTERVAL '90 days'`
	_, err := rm.db.Exec(query)
	return err
}
