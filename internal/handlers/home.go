package handlers

import (
	"net/http"

	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/middleware"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userID := middleware.SessionManager.GetInt(r.Context(), "userID")
	if userID != 0 {
		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		return
	}
	components.HomeLanding().Render(r.Context(), w)
}
