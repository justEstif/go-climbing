package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
)

var SessionManager *scs.SessionManager

// InitSessionManager initializes the SCS session manager with PostgreSQL storage
func InitSessionManager(dbPool *pgxpool.Pool) {
	SessionManager = scs.New()
	// Use web_sessions table to avoid conflict with climbing sessions table
	SessionManager.Store = pgxstore.NewWithConfig(dbPool, pgxstore.Config{
		TableName:       "web_sessions",
		CleanUpInterval: 5 * time.Minute,
	})
	SessionManager.Lifetime = 24 * time.Hour * 7 // 7 days
	SessionManager.Cookie.Name = "session_id"
	SessionManager.Cookie.HttpOnly = true
	SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	SessionManager.Cookie.Secure = os.Getenv("ENV") == "production"
	SessionManager.Cookie.Persist = true
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSignedIn := SessionManager.GetInt(r.Context(), "userID") != 0
		ctx := context.WithValue(r.Context(), components.IsSignedInKey, isSignedIn)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireNoAuth redirects authenticated users away from guest-only pages (e.g. login, signup)
func RequireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if SessionManager.GetInt(r.Context(), "userID") != 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireOnboarding redirects already-onboarded users away from /onboarding.
// Must be placed after RequireAuth in the middleware chain.
func RequireOnboarding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := SessionManager.GetInt(r.Context(), "userID")
		user, err := database.DB.GetUser(r.Context(), int32(userID))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if user.SessionsPerWeek != 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth is a middleware that ensures the user is authenticated
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := SessionManager.GetInt(r.Context(), "userID")
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
