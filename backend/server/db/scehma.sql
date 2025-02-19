-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create security_logs table
CREATE TABLE IF NOT EXISTS security_logs (
    log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type TEXT NOT NULL CHECK (type IN ('low', 'medium', 'high', 'critical')),
    hostname TEXT NOT NULL,
    host_ip TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    cve_id TEXT,
    cve_description TEXT,
    incident_report TEXT NOT NULL,
    threat_name TEXT NOT NULL,
    threat_category TEXT NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE NOT NULL,
    remedied_at TIMESTAMP WITH TIME ZONE,
    incident_created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    severity TEXT NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    threat_action TEXT,
    threat_status TEXT CHECK (threat_status IN ('handled', 'unresolved', 'under investigation')),
    detection_source TEXT,
    pid INTEGER,
    process_name TEXT,
    danger_score INTEGER NOT NULL CHECK (danger_score BETWEEN 0 AND 100)
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_timestamp ON security_logs (timestamp);
CREATE INDEX IF NOT EXISTS idx_severity ON security_logs (severity);
CREATE INDEX IF NOT EXISTS idx_danger_score ON security_logs (danger_score);
