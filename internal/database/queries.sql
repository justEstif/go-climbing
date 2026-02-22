-- User Queries

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, current_max_grade, goal_grade, sessions_per_week, weaknesses)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET email = $1, password_hash = $2, current_max_grade = $3, goal_grade = $4, sessions_per_week = $5, weaknesses = $6
WHERE id = $7;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserProfile :exec
UPDATE users
SET current_max_grade = $1, goal_grade = $2, sessions_per_week = $3, weaknesses = $4
WHERE id = $5;

-- Session Queries

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: ListSessionsByUser :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY date DESC, session_number ASC;

-- name: CreateSession :one
INSERT INTO sessions (user_id, session_number, date, focus_type, planned_warmup, planned_main, planned_project)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateSession :exec
UPDATE sessions
SET session_number = $1, date = $2, focus_type = $3, planned_warmup = $4, planned_main = $5, planned_project = $6
WHERE id = $7;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: GetLatestSessionByUser :one
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY date DESC, session_number DESC
LIMIT 1;

-- Session Log Queries

-- name: GetSessionLog :one
SELECT * FROM session_logs
WHERE id = $1 LIMIT 1;

-- name: ListSessionLogsByUser :many
SELECT * FROM session_logs
WHERE user_id = $1
ORDER BY logged_at DESC;

-- name: ListSessionLogsBySession :many
SELECT * FROM session_logs
WHERE session_id = $1
ORDER BY logged_at DESC;

-- name: CreateSessionLog :one
INSERT INTO session_logs (session_id, user_id, routes_logged, new_max_grade, energy_level, skin_condition, soreness, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateSessionLog :exec
UPDATE session_logs
SET routes_logged = $1, new_max_grade = $2, energy_level = $3, skin_condition = $4, soreness = $5, notes = $6
WHERE id = $7;

-- name: DeleteSessionLog :exec
DELETE FROM session_logs
WHERE id = $1;

-- Learn Content Queries

-- name: GetLearnContent :one
SELECT * FROM learn_content
WHERE id = $1 LIMIT 1;

-- name: ListLearnContentByCategory :many
SELECT * FROM learn_content
WHERE category = $1
ORDER BY sort_order ASC, title ASC;

-- name: ListAllLearnContent :many
SELECT * FROM learn_content
ORDER BY category ASC, sort_order ASC, title ASC;

-- name: CreateLearnContent :one
INSERT INTO learn_content (category, title, content, video_url, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateLearnContent :exec
UPDATE learn_content
SET category = $1, title = $2, content = $3, video_url = $4, sort_order = $5
WHERE id = $6;

-- name: DeleteLearnContent :exec
DELETE FROM learn_content
WHERE id = $1;

-- Feedback Queries

-- name: CreateFeedback :exec
INSERT INTO feedback (user_id, message) VALUES ($1, $2);
