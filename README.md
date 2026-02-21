# Go Web Template

Modern Go web application template with Chi, Templ, missing.css, and HTMX.

## Features

- **Chi** - Lightweight, idiomatic HTTP router
- **Templ** - Type-safe Go templating
- **missing.css** - Classless CSS framework
- **HTMX** - Dynamic interactions without heavy JavaScript
- **PostgreSQL** - Robust relational database
- **sqlc** - Compile-time type-safe SQL
- **golang-migrate** - Database migrations
- **mise** - Unified tool and task management

## Prerequisites

- [mise](https://mise.jdx.dev/) installed
- [Docker](https://www.docker.com/) for PostgreSQL

## Quick Start

### 1. Clone with gonew

```bash
gonew github.com/justestif/go-web-template github.com/yourname/myproject
cd myproject
```

### 2. Install Tools

```bash
mise install
```

This installs:
- Go (latest)
- Bun (latest)
- templ (latest)
- sqlc (latest)
- golang-migrate (latest)
- air (latest)

### 3. Start PostgreSQL

```bash
docker-compose up -d
```

### 4. Setup Project

```bash
mise run setup
```

This will:
- Install Node dependencies
- Run database migrations
- Generate templ components
- Generate type-safe SQL code

### 5. Start Development

Open two terminal windows:

**Terminal 1 - Watch for file changes (optional):**
```bash
mise run templ
```

**Terminal 2 - Go Server:**
```bash
mise run dev
```

### 6. Visit Application

Open [http://localhost:3000](http://localhost:3000)

## Available Tasks

Run `mise tasks` to see all available tasks:

- `mise run dev` - Start development server with live reload
- `mise run templ` - Generate templ files
- `mise run db-migrate` - Run database migrations
- `mise run db-rollback` - Rollback last migration
- `mise run sqlc` - Generate type-safe SQL code
- `mise run setup` - Complete project setup
- `mise run build` - Build production binary

## Project Structure

```
.
├── cmd/web/              # Application entry point
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   └── database/        # Database queries & connection
├── components/          # Templ templates
├── migrations/          # Database migrations
├── mise.toml           # Tool & task configuration
└── docker-compose.yml  # PostgreSQL setup
```

## Development Notes

### missing.css

Forms use Gorilla CSRF middleware:
- Token field name: `gorilla.csrf.Token`
- Access in templates: `csrf.Token(r)`
- Automatically validated on POST/PUT/DELETE
- Set `secure=true` in production (HTTPS only)

### Database

Sample migration creates a `users` table. See `internal/database/queries.sql` for example queries.

## Production Deployment

1. Set environment variables:
   ```bash
   export DATABASE_URL="postgres://..."
   export CSRF_KEY="your-32-byte-secret-key"
   export PORT="8080"
   ```

2. Update CSRF middleware to use `secure=true` in `cmd/web/main.go`

3. Build production binary:
   ```bash
   mise run build
   ```

4. Run migrations:
   ```bash
   mise run db-migrate
   ```

5. Start server:
   ```bash
   ./bin/app
   ```

## License

MIT
