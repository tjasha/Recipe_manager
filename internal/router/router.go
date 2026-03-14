package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tjasha/Recipe_manager/internal/handler"
)

func New(app *handler.Application) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// Initialize handlers
	h := handler.NewHandler(app)

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})

		// Authentication routes
		api.Post("/auth/google/verify", h.VerifyGoogleToken)
		api.Post("/auth/logout", h.Logout)

	})

	return r
}
