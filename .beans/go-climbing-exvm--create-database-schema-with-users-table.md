---
# go-climbing-exvm
title: Create database schema with users table
status: todo
type: task
priority: normal
created_at: 2026-02-21T23:00:29Z
updated_at: 2026-02-21T23:02:23Z
parent: go-climbing-rbkl
---

Set up PostgreSQL schema with users table including fields: id, email, password_hash, created_at, updated_at. Create migration file and run it.

```sql
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

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    session_number INTEGER NOT NULL,
    date DATE NOT NULL,
    focus_type VARCHAR(50) NOT NULL,
    planned_warmup JSONB NOT NULL,
    planned_main JSONB NOT NULL,
    planned_project JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE session_logs (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id),
    user_id INTEGER REFERENCES users(id),
    routes_logged JSONB NOT NULL,
    new_max_grade INTEGER,
    energy_level INTEGER CHECK (energy_level BETWEEN 1 AND 5),
    skin_condition VARCHAR(50),
    soreness INTEGER CHECK (soreness BETWEEN 1 AND 5),
    notes TEXT,
    logged_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE learn_content (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL, -- 'training_types', 'hold_types', 'techniques'
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    video_url VARCHAR(500), -- YouTube/external link
    sort_order INTEGER DEFAULT 0
);
```
