package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var SessionManager *scs.SessionManager

// InitSessionManager initializes the SCS session manager with PostgreSQL storage
func InitSessionManager(dbPool *pgxpool.Pool) {
	SessionManager = scs.New()
	SessionManager.Store = pgxstore.New(dbPool)
	SessionManager.Lifetime = 24 * time.Hour * 7 // 7 days
	SessionManager.Cookie.Name = "session_id"
	SessionManager.Cookie.HttpOnly = true
	SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	SessionManager.Cookie.Secure = os.Getenv("ENV") == "production"
	SessionManager.Cookie.Persist = true
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
