---
# go-climbing-fl6t
title: Add feedback form page
status: completed
type: feature
priority: normal
created_at: 2026-02-22T20:19:55Z
updated_at: 2026-02-22T20:22:36Z
---

Implement /feedback page for signed-in users: migration, SQL query, templ components, handler, router, navbar link

## Summary of Changes

- Created  and 
- Appended  query to ; regenerated sqlc
- Created  with , ,  components
- Created  with  (GET) and  (POST)
- Registered  GET and POST routes in  under 
- Added Feedback nav link to  inside the signed-in block
