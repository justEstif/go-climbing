---
# go-climbing-8i7n
title: logging
status: completed
type: task
priority: normal
created_at: 2026-02-21T22:54:50Z
updated_at: 2026-02-22T01:58:18Z
blocked_by:
    - go-climbing-cxyd
---

- [x] Build log form (steppers, sliders)
- [x] Associate logs with user_id
- [x] Test: log session, check DB shows correct user

## Summary of Changes

- /sessions lists all user sessions grouped into Today/Past with Logged chips linking to /sessions/{id}
- /sessions/{id} shows the session plan and log summary (energy, soreness, skin, notes) with Log/Update CTA
- /sessions/log fetches by session_id query param and pre-fills form when an existing log is found
- LogSubmit creates or updates session_log via log_id hidden field, redirects to /sessions/{id}
- Added energyValue, sorenessValue, skinSelected helpers for pre-filling the log form
