---
# go-climbing-zgqp
title: Create session middleware
status: completed
type: task
priority: normal
created_at: 2026-02-21T23:00:30Z
updated_at: 2026-02-21T23:46:49Z
parent: go-climbing-rbkl
---

Build session middleware to protect authenticated routes. Handle session creation on login, validation on protected routes, and cleanup on logout.



## Summary of Changes

Implemented SCS session middleware with PostgreSQL-backed sessions:

1. **internal/middleware/session.go** (NEW)
   - Created session manager initialization with 
   - Configured PostgreSQL pgxstore for server-side session storage
   - Set session lifetime to 7 days with secure cookie settings
   - Implemented  middleware to protect authenticated routes

2. **internal/handlers/logout.go** (NEW)
   - Created logout handler that destroys session using 
   - Redirects to home page after logout

3. **cmd/web/main.go**
   - Added  call after database initialization
   - Added  middleware for automatic session handling
   - Added  POST route

4. **Dependencies added to go.mod**
   - 
   - 

Key features:
- Server-side session storage in PostgreSQL (sessions persist across restarts)
- Automatic session token renewal on login (prevents session fixation)
- Secure cookie settings (HttpOnly, SameSite Strict, Secure in production)
