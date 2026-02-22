package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func SignupForm(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	components.Signup(csrfToken).Render(r.Context(), w)
}

func SignupSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.SignupError("Invalid form data").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if email == "" || password == "" {
		components.SignupError("Email and password are required").Render(r.Context(), w)
		return
	}

	if password != confirmPassword {
		components.SignupError("Passwords do not match").Render(r.Context(), w)
		return
	}

	if len(password) < 8 {
		components.SignupError("Password must be at least 8 characters").Render(r.Context(), w)
		return
	}

	// Check if user already exists
	existingUser, err := database.DB.GetUserByEmail(r.Context(), email)
	if err == nil && existingUser.ID != 0 {
		components.SignupError("Email already registered").Render(r.Context(), w)
		return
	}

	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		components.SignupError("Error processing password").Render(r.Context(), w)
		return
	}

	// Create user in database
	user, err := database.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:           email,
		PasswordHash:    string(hashedPassword),
		CurrentMaxGrade: 0,
		GoalGrade:       0,
		SessionsPerWeek: 0,
		Weaknesses:      []byte("[]"),
	})
	if err != nil {
		components.SignupError("Error creating account").Render(r.Context(), w)
		return
	}

	// Store user ID in session to automatically log them in
	middleware.SessionManager.Put(r.Context(), "userID", int(user.ID))

	w.Header().Set("HX-Redirect", "/onboarding")
	w.WriteHeader(http.StatusOK)
}
