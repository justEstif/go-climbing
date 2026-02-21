-- Migration 002: Create comprehensive climbing app schema
-- Drops the simple users table from migration 001 and recreates with full schema

-- Drop existing simple users table and related objects
DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;

-- Create comprehensive users table with climbing profile
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    current_max_grade INTEGER NOT NULL,
    goal_grade INTEGER NOT NULL,
    sessions_per_week INTEGER NOT NULL,
    weaknesses JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create training sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    session_number INTEGER NOT NULL,
    date DATE NOT NULL,
    focus_type VARCHAR(50) NOT NULL,
    planned_warmup JSONB NOT NULL,
    planned_main JSONB NOT NULL,
    planned_project JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for sessions
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_date ON sessions(date);

-- Create completed session logs table
CREATE TABLE session_logs (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    routes_logged JSONB NOT NULL,
    new_max_grade INTEGER,
    energy_level INTEGER CHECK (energy_level BETWEEN 1 AND 5),
    skin_condition VARCHAR(50),
    soreness INTEGER CHECK (soreness BETWEEN 1 AND 5),
    notes TEXT,
    logged_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for session_logs
CREATE INDEX idx_session_logs_session_id ON session_logs(session_id);
CREATE INDEX idx_session_logs_user_id ON session_logs(user_id);
CREATE INDEX idx_session_logs_logged_at ON session_logs(logged_at);

-- Create educational content table
CREATE TABLE learn_content (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    video_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0
);

-- Create index for learn_content
CREATE INDEX idx_learn_content_category ON learn_content(category);
CREATE INDEX idx_learn_content_sort_order ON learn_content(sort_order);
