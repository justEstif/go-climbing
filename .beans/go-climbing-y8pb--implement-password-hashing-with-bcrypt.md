---
# go-climbing-y8pb
title: Implement password hashing with bcrypt
status: completed
type: task
priority: normal
created_at: 2026-02-21T23:00:30Z
updated_at: 2026-02-21T23:46:45Z
parent: go-climbing-rbkl
---

Add bcrypt password hashing for secure password storage. Implement hashing on signup and verification on login.



## Summary of Changes

Implemented bcrypt password hashing:

1. **internal/handlers/signup.go**
   - Added bcrypt import for password hashing
   - Added check for existing users with 
   - Implemented  to hash passwords with default cost (10)
   - Create user in database with hashed password using 
   - Store user ID in session after signup for auto-login

2. **internal/handlers/login.go**
   - Added bcrypt import for password verification
   - Implemented  to verify passwords
   - Added session token renewal to prevent session fixation attacks
   - Store user ID in session after successful login
   - Removed placeholder authentication code

Dependencies:
-  (already available as transitive dependency)
