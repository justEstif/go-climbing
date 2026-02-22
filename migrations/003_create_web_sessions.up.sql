-- Migration: Create SCS web sessions table
-- This table stores server-side session data for the SCS session manager

CREATE TABLE IF NOT EXISTS web_sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

-- Index for efficient cleanup of expired sessions
CREATE INDEX IF NOT EXISTS idx_web_sessions_expiry ON web_sessions(expiry);
