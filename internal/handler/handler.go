package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5"
	"github.com/tjasha/Recipe_manager/internal/config"
	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"google.golang.org/api/idtoken"
)

// Application contains all dependencies for the application.
type Application struct {
	Config  *config.Config
	Session *scs.SessionManager
	DB      repository.Repository
}

// Handler allows access to Application instance in handler functions.
type Handler struct {
	App *Application
}

// New creates a new Handler.
func New(app *Application) *Handler {
	return &Handler{
		App: app,
	}
}

// Struct for receiving google token from the frontend
type GoogleToken struct {
	Credential string `json:"credential"`
}

// UserResponse is the data sent back to the frontend after successful login
type UserResponse struct {
	ID          uint   `json:"id"`
	UserName    string `json:"name"`
	Email       string `json:"email"`
	AccessLevel int    `json:"accessLevel"`
}

// VerifyGoogleToken handles the verification of the Google token
func (h *Handler) VerifyGoogleToken(w http.ResponseWriter, r *http.Request) {
	var token GoogleToken

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &token)
	if err != nil {
		log.Printf("ERROR: Cannot unmarshal token: %v", err)
		http.Error(w, "Cannot read token", http.StatusBadRequest)
		return
	}

	// Validate the token
	payload, err := idtoken.Validate(r.Context(), token.Credential, h.App.Config.GoogleOauthClientID)
	if err != nil {
		log.Printf("ERROR: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := payload.Claims
	log.Println(claims)

	// User logic
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		log.Printf("ERROR: Email not found or not a string in token claims")
		http.Error(w, "Email not found in token", http.StatusBadRequest)
		return
	}
	name, ok := claims["name"].(string)
	if !ok {
		name = "New User" // Default name if it doesn't exist
	}

	user, err := h.App.DB.GetUserByEmail(r.Context(), email)

	if err != nil {
		// If user does not exist, create a new one
		if errors.Is(err, pgx.ErrNoRows) {
			googleID := payload.Subject
			newUser := &model.User{
				UserName:    name,
				Email:       email,
				GoogleID:    &googleID,
				AccessLevel: 1, // Default access level for new users - chef
			}
			newID, createErr := h.App.DB.CreateUser(r.Context(), newUser)
			if createErr != nil {
				log.Printf("ERROR: Failed to create user: %v", createErr)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
			newUser.ID = uint(newID)
			user = newUser
		} else {
			// For any other database error
			log.Printf("ERROR: Database error on GetUserByEmail: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Renew token to prevent session fixation attacks.
	err = h.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to renew session token: %v", err)
		http.Error(w, "Failed to renew session token", http.StatusInternalServerError)
		return
	}

	// Store user info in the session
	h.App.Session.Put(r.Context(), "userID", user.ID)
	h.App.Session.Put(r.Context(), "accessLevel", user.AccessLevel)

	// Respond to Frontend with user info
	response := UserResponse{
		ID:          user.ID,
		UserName:    user.UserName,
		Email:       user.Email,
		AccessLevel: user.AccessLevel,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout destroys the user's session.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Destroy the session data
	err := h.App.Session.Destroy(r.Context())
	if err != nil {
		http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
		return
	}

	// Renew the token to ensure the old session is completely invalidated.
	err = h.App.Session.RenewToken(r.Context())
	if err != nil {
		http.Error(w, "Failed to renew token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logout successful"}`))
}

// ShowAllRecipes shows all recipes on the homepage.
func (h *Handler) ShowAllRecipes(w http.ResponseWriter, r *http.Request) {

	recipes, err := h.App.DB.GetAllRecipes(r.Context())
	log.Println(recipes)
	if err != nil {
		http.Error(w, "Failed to retrieve recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
	log.Println(recipes)
}

func (h *Handler) ShowFullRecipe(w http.ResponseWriter, r *http.Request, id int64) {

	recipe, err := h.App.DB.GetRecipeByID(id)
	log.Println(recipe)
	if err != nil {
		http.Error(w, "Failed to retrieve recipe", http.StatusInternalServerError)
		return
	}

	ingredients, err := h.App.DB.GetIngredientsByRecipeID(r.Context(), id)
	log.Println(ingredients)
	if err != nil {
		http.Error(w, "Failed to retrieve ingredients in the recipe", http.StatusInternalServerError)
		return
	}

	recipe.Ingredients = ingredients

	instructions, err := h.App.DB.GetInstructionsByRecipeID(r.Context(), id)
	log.Println(instructions)
	if err != nil {
		http.Error(w, "Failed to retrieve instructions", http.StatusInternalServerError)
		return
	}

	recipe.Instructions = instructions

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
	log.Println(recipe)
}
