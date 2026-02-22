package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	climbsession "github.com/justestif/go-climbing/internal/session"
)

func LogForm(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.URL.Query().Get("session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil || sessionID <= 0 {
		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		return
	}

	session, err := database.DB.GetSession(r.Context(), int32(sessionID))
	if err != nil {
		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		return
	}

	logs, err := database.DB.ListSessionLogsBySession(r.Context(), pgtype.Int4{Int32: int32(sessionID), Valid: true})
	var existingLog *database.SessionLog
	if err == nil && len(logs) > 0 {
		existingLog = &logs[0]
	}

	csrfToken := csrf.Token(r)
	components.LogForm(csrfToken, session, existingLog).Render(r.Context(), w)
}

func LogSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.LogError("Invalid form data").Render(r.Context(), w)
		return
	}

	sessionIDStr := r.FormValue("session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil || sessionID <= 0 {
		components.LogError("Invalid session ID").Render(r.Context(), w)
		return
	}

	// Parse multi-value route fields, filtering rows with empty grade
	grades := r.Form["route_grade"]
	styles := r.Form["route_style"]
	counts := r.Form["route_count"]

	var routes []climbsession.LoggedRoute
	for i, g := range grades {
		if g == "" {
			continue
		}
		grade, err := strconv.Atoi(g)
		if err != nil {
			continue
		}
		style := ""
		if i < len(styles) {
			style = styles[i]
		}
		count := 1
		if i < len(counts) {
			if c, err := strconv.Atoi(counts[i]); err == nil && c > 0 {
				count = c
			}
		}
		routes = append(routes, climbsession.LoggedRoute{Grade: grade, Style: style, Count: count})
	}

	routesJSON, err := climbsession.EncodeRoutesLogged(routes)
	if err != nil {
		components.LogError("Error encoding routes").Render(r.Context(), w)
		return
	}

	energyLevel, _ := strconv.Atoi(r.FormValue("energy_level"))
	soreness, _ := strconv.Atoi(r.FormValue("soreness"))
	skinCondition := r.FormValue("skin_condition")
	notes := r.FormValue("notes")

	newMaxGradeStr := r.FormValue("new_max_grade")
	var newMaxGrade pgtype.Int4
	if newMaxGradeStr != "" {
		if g, err := strconv.Atoi(newMaxGradeStr); err == nil {
			newMaxGrade = pgtype.Int4{Int32: int32(g), Valid: true}
		}
	}

	logIDStr := r.FormValue("log_id")
	if logIDStr != "" {
		logID, err := strconv.Atoi(logIDStr)
		if err != nil {
			components.LogError("Invalid log ID").Render(r.Context(), w)
			return
		}
		err = database.DB.UpdateSessionLog(r.Context(), database.UpdateSessionLogParams{
			ID:            int32(logID),
			RoutesLogged:  routesJSON,
			NewMaxGrade:   newMaxGrade,
			EnergyLevel:   pgtype.Int4{Int32: int32(energyLevel), Valid: energyLevel > 0},
			SkinCondition: pgtype.Text{String: skinCondition, Valid: skinCondition != ""},
			Soreness:      pgtype.Int4{Int32: int32(soreness), Valid: soreness > 0},
			Notes:         pgtype.Text{String: notes, Valid: notes != ""},
		})
		if err != nil {
			components.LogError("Error updating session log").Render(r.Context(), w)
			return
		}
	} else {
		userIDInt := int32(middleware.SessionManager.GetInt(r.Context(), "userID"))
		_, err = database.DB.CreateSessionLog(r.Context(), database.CreateSessionLogParams{
			SessionID:     pgtype.Int4{Int32: int32(sessionID), Valid: true},
			UserID:        pgtype.Int4{Int32: userIDInt, Valid: true},
			RoutesLogged:  routesJSON,
			NewMaxGrade:   newMaxGrade,
			EnergyLevel:   pgtype.Int4{Int32: int32(energyLevel), Valid: energyLevel > 0},
			SkinCondition: pgtype.Text{String: skinCondition, Valid: skinCondition != ""},
			Soreness:      pgtype.Int4{Int32: int32(soreness), Valid: soreness > 0},
			Notes:         pgtype.Text{String: notes, Valid: notes != ""},
		})
		if err != nil {
			components.LogError("Error saving session log").Render(r.Context(), w)
			return
		}

		// Generate next session if this is the latest one (best-effort)
		latestSession, latestErr := database.DB.GetLatestSessionByUser(r.Context(), pgtype.Int4{Int32: userIDInt, Valid: true})
		if latestErr == nil && latestSession.ID == int32(sessionID) {
			currentSession, sessErr := database.DB.GetSession(r.Context(), int32(sessionID))
			user, userErr := database.DB.GetUser(r.Context(), userIDInt)
			if sessErr == nil && userErr == nil {
				weaknesses, _ := climbsession.ParseWeaknesses(user.Weaknesses)
				params, genErr := climbsession.GenerateNextSession(
					userIDInt,
					currentSession.SessionNumber,
					user.CurrentMaxGrade,
					user.GoalGrade,
					weaknesses,
					currentSession.Date.Time,
					user.SessionsPerWeek,
				)
				if genErr == nil {
					database.DB.CreateSession(r.Context(), params) //nolint:errcheck // best-effort
				}
			}
		}
	}

	w.Header().Set("HX-Redirect", "/sessions/"+sessionIDStr)
	w.WriteHeader(http.StatusOK)
}
