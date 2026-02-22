---
# go-climbing-64fe
title: progress + learn
status: completed
type: task
priority: normal
created_at: 2026-02-21T22:50:59Z
updated_at: 2026-02-22T02:12:15Z
blocked_by:
    - go-climbing-8i7n
---

- [x] Build progress page (user's personal data only)
- [x] Add Chart.js visualization
- [x] Build Learn section structure
- [x] Seed DB with 5-10 learn items (training types, holds, techniques)
- [x] Add YouTube links for technique videos

## Summary of Changes

- Added `/progress` page: session count, current/goal grade overview, grade progression line chart (Chart.js, V0–V17 y-axis), and energy/soreness wellness chart with spanGaps
- Added `/learn` list page: items grouped by category with Video chips for entries with YouTube links
- Added `/learn/{id}` detail page: renders developer-authored HTML content and responsive YouTube iframe embeds
- Created migration 004 seeding 10 learn items across 5 categories (Footwork, Body Positioning, Training Methods, Hold Types, Mental Skills)
- Updated navbar to show Sessions / Progress / Learn links for signed-in users
- Registered `/progress`, `/learn`, `/learn/{id}` routes in the RequireAuth group
