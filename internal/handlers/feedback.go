package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
)

func FeedbackForm(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	components.Feedback(csrfToken).Render(r.Context(), w)
}

func FeedbackSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.FeedbackError("Invalid form data").Render(r.Context(), w)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		components.FeedbackError("Message cannot be empty").Render(r.Context(), w)
		return
	}

	userID := middleware.SessionManager.GetInt(r.Context(), "userID")
	err := database.DB.CreateFeedback(r.Context(), database.CreateFeedbackParams{
		UserID:  pgtype.Int4{Int32: int32(userID), Valid: true},
		Message: message,
	})
	if err != nil {
		components.FeedbackError("Failed to submit feedback. Please try again.").Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Retarget", "#feedback-form")
	w.Header().Set("HX-Reswap", "outerHTML")
	components.FeedbackSuccess().Render(r.Context(), w)
}
