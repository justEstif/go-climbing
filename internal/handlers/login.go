package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func LoginForm(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	components.Login(csrfToken).Render(r.Context(), w)
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		components.LoginError("Invalid form data").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		components.LoginError("Email and password are required").Render(r.Context(), w)
		return
	}

	// Fetch user by email
	user, err := database.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		// Don't reveal whether email exists or not for security
		components.LoginError("Invalid email or password").Render(r.Context(), w)
		return
	}

	// Verify password with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		components.LoginError("Invalid email or password").Render(r.Context(), w)
		return
	}

	// Renew session token to prevent session fixation
	if err := middleware.SessionManager.RenewToken(r.Context()); err != nil {
		components.LoginError("Error creating session").Render(r.Context(), w)
		return
	}

	// Store user ID in session
	middleware.SessionManager.Put(r.Context(), "userID", int(user.ID))

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
