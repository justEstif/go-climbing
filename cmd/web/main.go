package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/handlers"
	customMiddleware "github.com/justestif/go-climbing/internal/middleware"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize session manager
	customMiddleware.InitSessionManager(database.Pool)

	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Session middleware - handles loading and saving session data
	r.Use(customMiddleware.SessionManager.LoadAndSave)

	// Auth middleware - sets isSignedIn in context for all requests
	r.Use(customMiddleware.AuthMiddleware)

	// CSRF protection - set secure=true in production
	csrfKey := []byte(os.Getenv("CSRF_KEY"))
	if len(csrfKey) != 32 {
		log.Fatal("CSRF_KEY must be exactly 32 bytes long")
	}
	csrfMw := customMiddleware.SetupCSRF(csrfKey, false)

	// Static files (no CSRF needed)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes with CSRF
	r.Group(func(r chi.Router) {
		r.Use(csrfMw)
		r.Get("/", handlers.Home)
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RequireNoAuth)
			r.Get("/signup", handlers.SignupForm)
			r.Post("/signup", handlers.SignupSubmit)
			r.Get("/login", handlers.LoginForm)
			r.Post("/login", handlers.LoginSubmit)
		})
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RequireAuth)
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.RequireOnboarding)
				r.Get("/onboarding", handlers.OnboardingForm)
				r.Post("/onboarding", handlers.OnboardingSubmit)
			})
			r.Get("/sessions", handlers.SessionsPage)
			r.Get("/sessions/log", handlers.LogForm)
			r.Post("/sessions/log", handlers.LogSubmit)
			r.Get("/sessions/{id}", handlers.SessionDetail)
			r.Get("/progress", handlers.ProgressPage)
			r.Get("/learn", handlers.LearnListPage)
			r.Get("/learn/{id}", handlers.LearnDetailPage)
			r.Get("/feedback", handlers.FeedbackForm)
			r.Post("/feedback", handlers.FeedbackSubmit)
		})
	})

	// Logout route (no CSRF needed for logout)
	r.Post("/logout", handlers.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
