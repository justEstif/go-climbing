# Go Climbing

A climbing training app to help climbers track progress, plan sessions, and level up their skills.

**Live:** http://go-climbing.estifanos.cc

## What It Does

- **Onboarding** — set your current max grade, goal grade, sessions per week, and weaknesses
- **Session planning** — AI-generated training plans based on your profile
- **Session logging** — log completed sessions and track what you actually did
- **Progress tracking** — visualize improvement over time
- **Learn** — educational content on technique, training concepts, and injury prevention

## Built With

- **[Go](https://go.dev/)** — backend language
- **[Chi](https://github.com/go-chi/chi)** — HTTP router
- **[Templ](https://templ.guide/)** — type-safe HTML templating
- **[HTMX](https://htmx.org/)** — dynamic interactions without heavy JavaScript
- **[missing.css](https://missing.style/)** — classless CSS framework
- **[PostgreSQL](https://www.postgresql.org/)** — database
- **[sqlc](https://sqlc.dev/)** — compile-time type-safe SQL queries
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — database migrations
- **[SCS](https://github.com/alexedwards/scs)** — session management
- **[mise](https://mise.jdx.dev/)** — tool and task management
- **[Dokku](https://dokku.com/)** — self-hosted deployment

## Local Development

**Prerequisites:** [mise](https://mise.jdx.dev/), [Docker](https://www.docker.com/)

```bash
# Install tools
mise install

# Start PostgreSQL
docker-compose up -d

# Run migrations, generate templ + sqlc
mise run setup

# Start dev server with live reload
mise run dev
```

Open [http://localhost:3000](http://localhost:3000)

## Available Tasks

```bash
mise run dev          # Start dev server with live reload
mise run build        # Build production binary to bin/app
mise run templ        # Regenerate templ components (after editing .templ files)
mise run sqlc         # Regenerate type-safe SQL (after editing queries.sql)
mise run db-migrate   # Apply pending migrations
mise run db-rollback  # Roll back last migration
mise run setup        # Full initial setup (migrations + templ + sqlc)
```

## Project Structure

```
.
├── cmd/web/              # Entry point (router, middleware, server)
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/       # Auth, session, CSRF middleware
│   ├── database/         # sqlc-generated models and queries
│   └── session/          # Session planning logic
├── components/           # Templ templates
├── migrations/           # SQL migration files
├── mise.toml             # Tool & task configuration
└── docker-compose.yml    # Local PostgreSQL
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `CSRF_KEY` | Yes | 32-byte CSRF secret key |
| `PORT` | No | Server port (default: 3000) |
| `ENV` | No | Set to `production` for secure HTTPS-only cookies |
