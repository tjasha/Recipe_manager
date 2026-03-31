package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tjasha/Recipe_manager/internal/handler"
)

func New(app *handler.Application) http.Handler {
	r := chi.NewRouter()

	// --- Middleware ---
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// CORS middleware - allows different origin of frontend and backend
	r.Use(cors.Handler(cors.Options{
		// Allowed origin - frontend
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		// Allows cookies / sessions
		AllowCredentials: true,
		MaxAge:           500,
	}))

	// Session middleware
	r.Use(app.Session.LoadAndSave)

	// Initialize handlers
	h := handler.New(app)

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(`{"status":"ok"}`))
			if err != nil {
				log.Printf("health check write failed: %v", err)
				return
			}
		})

		// Authentication routes
		api.Post("/auth/google/verify", h.VerifyGoogleToken)
		api.Post("/auth/logout", h.Logout)

		// Recipes routes
		api.Get("/recipes", h.ShowAllPublishedRecipes)
		api.Get("/recipe/{id}", h.ShowFullRecipe)
		//chef routes
		api.Get("/myrecipes", h.ShowRecipesOfTheUser)
		api.Get("/myrecipe/{id}", h.ShowFullRecipe)
		api.Get("/ingredients", h.GetAllIngredients)
		api.Post("/createRecipe", h.SaveRecipe)
		api.Delete("/deleteRecipe/{id}", h.DeleteRecipe)
		api.Patch("/myrecipe/{id}/publish", h.UpdatePublishRecipe)
		api.Put("/myrecipe/{id}", h.EditRecipe)
		//admin

	})

	return r
}
