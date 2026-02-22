package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	climbsession "github.com/justestif/go-climbing/internal/session"
)

func OnboardingForm(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	components.Onboarding(csrfToken).Render(r.Context(), w)
}

func OnboardingSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.OnboardingError("Invalid form data").Render(r.Context(), w)
		return
	}

	currentGradeStr := r.FormValue("current_max_grade")
	goalGradeStr := r.FormValue("goal_grade")
	sessionsPerWeekStr := r.FormValue("sessions_per_week")
	weaknessValues := r.Form["weaknesses"]

	currentGrade, err := strconv.Atoi(currentGradeStr)
	if err != nil || currentGrade < 0 || currentGrade > 17 {
		components.OnboardingError("Invalid current grade").Render(r.Context(), w)
		return
	}

	goalGrade, err := strconv.Atoi(goalGradeStr)
	if err != nil || goalGrade < 0 || goalGrade > 17 {
		components.OnboardingError("Invalid goal grade").Render(r.Context(), w)
		return
	}

	if goalGrade < currentGrade {
		components.OnboardingError("Goal grade must be at least your current grade").Render(r.Context(), w)
		return
	}

	sessionsPerWeek, err := strconv.Atoi(sessionsPerWeekStr)
	if err != nil || sessionsPerWeek < 1 || sessionsPerWeek > 7 {
		components.OnboardingError("Sessions per week must be between 1 and 7").Render(r.Context(), w)
		return
	}

	userID := middleware.SessionManager.GetInt(r.Context(), "userID")

	weaknessesJSON, err := climbsession.MarshalWeaknesses(weaknessValues)
	if err != nil {
		components.OnboardingError("Error processing weaknesses").Render(r.Context(), w)
		return
	}

	err = database.DB.UpdateUserProfile(r.Context(), database.UpdateUserProfileParams{
		CurrentMaxGrade: int32(currentGrade),
		GoalGrade:       int32(goalGrade),
		SessionsPerWeek: int32(sessionsPerWeek),
		Weaknesses:      weaknessesJSON,
		ID:              int32(userID),
	})
	if err != nil {
		components.OnboardingError("Error saving profile").Render(r.Context(), w)
		return
	}

	params, err := climbsession.GenerateFirstSession(
		int32(userID),
		int32(currentGrade),
		int32(goalGrade),
		weaknessValues,
		time.Now(),
	)
	if err != nil {
		components.OnboardingError("Error generating session plan").Render(r.Context(), w)
		return
	}

	_, err = database.DB.CreateSession(r.Context(), params)
	if err != nil {
		components.OnboardingError("Error saving session plan").Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Redirect", "/sessions")
	w.WriteHeader(http.StatusOK)
}
