package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tjasha/Recipe_manager/internal/handler"
)

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--> INCOMING REQUEST: Method=[%s], Path=[%s]", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware for security headers
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allows comunication with google SSO
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		next.ServeHTTP(w, r)
	})
}

//func manualCorsMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Nastavimo glave, ki jih pričakuje brskalnik
//		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
//		w.Header().Set("Access-Control-Allow-Credentials", "true")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Credentials")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
//
//		// Če je to "preflight" OPTIONS zahteva, jo takoj zaključimo
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		// Za vse ostale zahteve, jih samo pošljemo naprej
//		next.ServeHTTP(w, r)
//	})
//}

func New(app *handler.Application) http.Handler {
	r := chi.NewRouter()

	// 1. KORAK: DODAMO DIAGNOSTIČNI LOGGER NA SAM VRH
	r.Use(requestLoggerMiddleware)

	//// ROCNI MIDDLEWARE
	//r.Use(manualCorsMiddleware)

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
		api.Post("/recipes/{id}/adjust-portion", h.AdjustServingSize)

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
