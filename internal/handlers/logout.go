package handlers

import (
	"net/http"

	"github.com/justestif/go-climbing/internal/middleware"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Destroy the session
	if err := middleware.SessionManager.Destroy(r.Context()); err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
