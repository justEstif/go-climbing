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
		r.Get("/about", handlers.About)
		r.Get("/contact", handlers.ContactForm)
		r.Post("/contact", handlers.ContactSubmit)
		r.Get("/signup", handlers.SignupForm)
		r.Post("/signup", handlers.SignupSubmit)
		r.Get("/login", handlers.LoginForm)
		r.Post("/login", handlers.LoginSubmit)
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
