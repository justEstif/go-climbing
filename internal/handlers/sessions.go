package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	climbsession "github.com/justestif/go-climbing/internal/session"
)

func SessionsPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.SessionManager.GetInt(r.Context(), "userID")

	user, err := database.DB.GetUser(r.Context(), int32(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if user.SessionsPerWeek == 0 {
		http.Redirect(w, r, "/onboarding", http.StatusSeeOther)
		return
	}

	sessions, err := database.DB.ListSessionsByUser(r.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		components.SessionsDashboard(user, nil, "", nil, nil).Render(r.Context(), w)
		return
	}

	userLogs, err := database.DB.ListSessionLogsByUser(r.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		userLogs = nil
	}

	loggedIDs := map[int32]bool{}
	for _, l := range userLogs {
		if l.SessionID.Valid {
			loggedIDs[l.SessionID.Int32] = true
		}
	}

	today := time.Now().Format("2006-01-02")

	// Find the earliest upcoming date (today or future).
	nextDate := ""
	for _, s := range sessions {
		d := s.Date.Time.Format("2006-01-02")
		if d >= today && (nextDate == "" || d < nextDate) {
			nextDate = d
		}
	}

	// Format for display, e.g. "Feb 22, 2026".
	nextDateLabel := ""
	if nextDate != "" {
		if t, err := time.Parse("2006-01-02", nextDate); err == nil {
			nextDateLabel = t.Format("Jan 2, 2006")
		}
	}

	var nextSessions, pastSessions []database.Session
	for _, s := range sessions {
		d := s.Date.Time.Format("2006-01-02")
		if d >= today {
			nextSessions = append(nextSessions, s)
		} else {
			pastSessions = append(pastSessions, s)
		}
	}

	components.SessionsDashboard(user, nextSessions, nextDateLabel, pastSessions, loggedIDs).Render(r.Context(), w)
}

func SessionDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	session, err := database.DB.GetSession(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	plan, planErr := climbsession.DecodeSessionPlan(session)
	var planPtr *climbsession.SessionPlan
	if planErr == nil {
		planPtr = &plan
	}

	logs, err := database.DB.ListSessionLogsBySession(r.Context(), pgtype.Int4{Int32: int32(id), Valid: true})
	var existingLog *database.SessionLog
	if err == nil && len(logs) > 0 {
		existingLog = &logs[0]
	}

	components.SessionDetail(session, planPtr, existingLog).Render(r.Context(), w)
}
