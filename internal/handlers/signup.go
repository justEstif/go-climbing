package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justestif/go-climbing/components"
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

	// Note: Actual user creation with bcrypt hashing will be implemented in bean go-climbing-y8pb
	// For now, return success to demonstrate the form flow
	components.SignupSuccess().Render(r.Context(), w)
}
