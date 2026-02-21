-- Migration 002 Down: Rollback climbing app schema

-- Drop tables in reverse order to handle foreign key constraints
DROP TABLE IF EXISTS learn_content CASCADE;
DROP TABLE IF EXISTS session_logs CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
