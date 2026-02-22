package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	climbsession "github.com/justestif/go-climbing/internal/session"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userID := middleware.SessionManager.GetInt(r.Context(), "userID")
	if userID == 0 {
		components.HomeLanding().Render(r.Context(), w)
		return
	}

	user, err := database.DB.GetUser(r.Context(), int32(userID))
	if err != nil {
		components.HomeLanding().Render(r.Context(), w)
		return
	}

	if user.SessionsPerWeek == 0 {
		http.Redirect(w, r, "/onboarding", http.StatusSeeOther)
		return
	}

	dbSession, err := database.DB.GetLatestSessionByUser(r.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		components.HomeDashboard(user, nil, nil).Render(r.Context(), w)
		return
	}

	plan, err := climbsession.DecodeSessionPlan(dbSession)
	if err != nil {
		components.HomeDashboard(user, &dbSession, nil).Render(r.Context(), w)
		return
	}

	components.HomeDashboard(user, &dbSession, &plan).Render(r.Context(), w)
}
