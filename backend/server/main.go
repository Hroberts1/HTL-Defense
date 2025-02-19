package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"
	"time"

	"HTL-Defense/backend/server/api"
	"HTL-Defense/backend/server/core/config"
	"HTL-Defense/backend/server/core/logging"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection using cfg.DatabaseURL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize retention manager
	retention := logging.NewRetentionManager(db, 90)
	retention.StartRetentionPolicy(ctx)

	// Initialize API handlers
	logsHandler := api.NewLogsHandler(db)

	// Set up routes
	http.HandleFunc("/api/logs", logsHandler.GetLogs)

	// Setup TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
