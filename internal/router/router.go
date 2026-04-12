package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tjasha/Recipe_manager/internal/handler"
)

// Middleware for security headers
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allows communication with google SSO
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		next.ServeHTTP(w, r)
	})
}

// New creates a new router.
func New(app *handler.Application) http.Handler {
	r := chi.NewRouter()

	// CORS middleware - allows different origin of frontend and backend
	r.Use(cors.Handler(cors.Options{
		// Allowed origin - frontend
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "Credentials", "X-CSRF-Token"},
		// Allows cookies / sessions
		AllowCredentials: true,
		MaxAge:           300,
		// to automatically answer OPTIONS requests
		OptionsPassthrough: false,
	}))

	// --- Middleware ---
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// use security headers
	r.Use(securityHeadersMiddleware)

	// Session middleware
	r.Use(app.Session.LoadAndSave)

	// Initialize handlers
	h := handler.New(app)

	r.Route("/api", func(api chi.Router) {

		// Accessible to all users
		// Health check
		api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"status":"ok"}`))
			if err != nil {
				log.Printf("health check write failed: %v", err)
				return
			}
		})
		// Authentication
		api.Post("/auth/google/verify", h.VerifyGoogleToken)
		api.Post("/auth/logout", h.Logout)
		api.Get("/auth/check-session", h.CheckSession)

		// Recipes routes
		api.Get("/recipes", h.ShowAllPublishedRecipes)
		api.Get("/recipe/{id}", h.ShowFullRecipe)
		api.Post("/recipes/{id}/adjust-portion", h.AdjustServingSize)
		api.Get("/ingredients", h.GetAllIngredients)

		// Routes for logged in users
		// Accessible to chefs (access level 1)

		api.Group(func(r chi.Router) {
			r.Use(h.RequireRole(1))

			r.Get("/myrecipes", h.ShowRecipesOfTheUser)
			r.Get("/myrecipe/{id}", h.ShowFullRecipe)
			r.Post("/createRecipe", h.SaveRecipe)
			r.Delete("/deleteRecipe/{id}", h.DeleteRecipe)
			r.Patch("/myrecipe/{id}/publish", h.UpdatePublishRecipe)
			r.Put("/myrecipe/{id}", h.EditRecipe)
		})

		// Accessible to admins (access level 0)
		api.Group(func(r chi.Router) {
			r.Use(h.RequireRole(0))

			r.Get("/admin/users", h.ReturnAllUsers)
			r.Delete("/admin/users/{id}", h.DeleteUser)
			r.Patch("/admin/users/{id}", h.UpdateUser)
		})
	})

	return r
}
