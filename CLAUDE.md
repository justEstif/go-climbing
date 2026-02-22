# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This project uses [Mise](https://mise.jdx.dev/) for task management (not Make).

```bash
mise run dev          # Start dev server with live reload (air)
mise run build        # Build production binary to bin/app
mise run setup        # Run migrations + generate templ + generate sqlc (initial setup)
mise run templ        # Regenerate templ components (run after editing .templ files)
mise run sqlc         # Regenerate type-safe SQL code (run after editing queries.sql)
mise run db-migrate   # Apply pending migrations
mise run db-rollback  # Roll back last migration
```

**Important:** After editing any `.templ` file, run `mise run templ` to regenerate the `_templ.go` files before building or running.

There are no tests in this codebase.

## Architecture

### Request Flow
```
HTTP Request
  → Chi router
  → Global middleware (Logger, Recoverer, RequestID, RealIP)
  → SCS SessionManager.LoadAndSave
  → AuthMiddleware (sets isSignedIn in context via components.IsSignedInKey)
  → [CSRF middleware on non-static routes]
  → Handler
  → Templ component render
```

### Routes
| Method | Path | Middleware | Handler |
|--------|------|------------|---------|
| GET | `/` | — | `Home` |
| GET/POST | `/signup` | RequireNoAuth | `SignupForm` / `SignupSubmit` |
| GET/POST | `/login` | RequireNoAuth | `LoginForm` / `LoginSubmit` |
| POST | `/logout` | — (no CSRF) | `Logout` |
| GET/POST | `/onboarding` | RequireAuth, RequireOnboarding | `OnboardingForm` / `OnboardingSubmit` |
| GET | `/sessions` | RequireAuth | `SessionsPage` |
| GET/POST | `/sessions/log` | RequireAuth | `LogForm` / `LogSubmit` |
| GET | `/sessions/{id}` | RequireAuth | `SessionDetail` |
| GET | `/progress` | RequireAuth | `ProgressPage` |
| GET | `/learn` | RequireAuth | `LearnListPage` |
| GET | `/learn/{id}` | RequireAuth | `LearnDetailPage` |
| GET/POST | `/feedback` | RequireAuth | `FeedbackForm` / `FeedbackSubmit` |

### Key Directories
- `cmd/web/main.go` — Entry point: router setup, middleware chain, server start
- `internal/handlers/` — HTTP handlers (home, signup, login, logout, onboarding, sessions, log, progress, learn, feedback)
- `internal/middleware/` — Session and CSRF middleware
- `internal/database/` — Database layer: `init.go` (connection), `models.go` and `queries.sql.go` (sqlc-generated)
- `internal/session/` — Session planner: generates structured climbing workouts (`SessionPlan` with `Warmup`, `Main`, `Project` phases) based on user grade, goal grade, and weaknesses; also handles deload logic and route logging encode/decode
- `components/` — Templ templates (`.templ` source + `_templ.go` generated) plus plain `.go` files: `context.go`, `helpers.go`, `learn_data.go`, `progress_data.go`
- `migrations/` — SQL migration files (golang-migrate format)

### Code Generation
- **sqlc**: Reads `internal/database/queries.sql` + `migrations/` schema → generates `models.go` and `queries.sql.go`. Edit the `.sql` file and run `mise run sqlc`.
- **templ**: `.templ` files compile to `_templ.go`. Never edit `_templ.go` directly.

### Database
PostgreSQL 16 via Docker (`docker-compose up -d`). Connection string: `DATABASE_URL` env var (default: `postgres://postgres:postgres@localhost:5432/appdb`).

Schema highlights:
- `users` — Auth + climbing profile (grade, goals, weaknesses as JSONB)
- `sessions` / `session_logs` — Training sessions with planned/actual workout data (JSONB columns)
- `learn_content` — Educational content categorized by type (seeded in migration 004)
- `feedback` — User feedback submissions
- `web_sessions` — SCS session storage (NOT the climbing `sessions` table)

### Auth System
- Sessions: SCS v2 with PostgreSQL backend (`web_sessions` table). Cookie: `session_id`, 7-day expiry, SameSite=Strict.
- Passwords: bcrypt at DefaultCost.
- CSRF: Gorilla CSRF, key from `CSRF_KEY` env var (32 bytes). Token passed to templates via `csrf.Token(r)`.
- `AuthMiddleware` sets `isSignedIn` via `components.IsSignedInKey` (typed context key) for all requests.
- `RequireAuth` redirects unauthenticated users to `/login` and stores `userID` in context using the plain string key `"userID"` (not the typed key).
- `RequireOnboarding` (applied after `RequireAuth`) redirects already-onboarded users away from `/onboarding`. Onboarding is complete when `users.sessions_per_week != 0`.

### Templating Conventions
- `components/context.go` provides `GetIsSignedIn(ctx)` helper for reading auth state in templates.
- Forms must include the CSRF hidden field — pass `csrf.Token(r)` from the handler to the templ component.
- Base layout is `components/layout.templ`; uses [missing.style](https://missing.style/) (classless CSS framework) + HTMX via CDN.
- HTMX redirect pattern: `w.Header().Set("HX-Redirect", "/path"); w.WriteHeader(http.StatusOK)`
- Climbing grades are Hueco V-grades stored as integers; use `gradeName(n)` → `"V0"`–`"V15"`.

## Environment Variables
| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `CSRF_KEY` | Yes | 32-byte CSRF secret key |
| `PORT` | No | Server port (default: 3000) |
| `ENV` | No | Set to `production` for secure HTTPS-only cookies |
