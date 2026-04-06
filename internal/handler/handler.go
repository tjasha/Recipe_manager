package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/tjasha/Recipe_manager/internal/config"
	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/repository"
	"github.com/tjasha/Recipe_manager/internal/service"
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

// CheckSession checks if the session is active and retrieve full user details to send back
func (h *Handler) CheckSession(w http.ResponseWriter, r *http.Request) {
	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	// If a session exists, retrieve full user details to send back
	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	// send back the data we have in the session
	username, ok := h.App.Session.Get(r.Context(), "username").(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}
	accessLevel, ok := h.App.Session.Get(r.Context(), "accessLevel").(int)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		ID:          uid,
		UserName:    username,
		AccessLevel: accessLevel,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error: fail to encode response: %v", err)
		return
	}
}

// RequireRole check access level of the user.
func (h *Handler) RequireRole(requiredLevel int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// check if user is logged in
			accessLevel := h.App.Session.Get(r.Context(), "accessLevel")
			if accessLevel == nil {
				http.Error(w, "Unauthorized: You must be logged in.", http.StatusUnauthorized)
				return
			}
			// check if user has the required access level
			userLevel, ok := accessLevel.(int)
			if !ok {
				http.Error(w, "Invalid session state.", http.StatusInternalServerError)
				return
			}
			// check if level is high enough
			if userLevel > requiredLevel {
				http.Error(w, "Forbidden: You do not have permission to perform this action.", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GoogleToken is a struct for receiving Google token from the frontend
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Cannot read body: %v", err)
		http.Error(w, "Cannot read request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &token)
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

	// Respond to Frontend with user info
	response := UserResponse{
		ID:          user.ID,
		UserName:    user.UserName,
		Email:       user.Email,
		AccessLevel: user.AccessLevel,
	}

	// Store user info in the session
	h.App.Session.Put(r.Context(), "userID", user.ID)
	h.App.Session.Put(r.Context(), "username", user.UserName)
	h.App.Session.Put(r.Context(), "accessLevel", user.AccessLevel)
	log.Println("Store user info in the session: ", h.App.Session.Get(r.Context(), "userID"), h.App.Session.Get(r.Context(), "username"), h.App.Session.Get(r.Context(), "accessLevel"))

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// Logout destroys the user's session.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("v logout handlerju")
	// destroy session
	err := h.App.Session.Destroy(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to destroy session: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ShowAllPublishedRecipes shows all published recipes on the homepage.
func (h *Handler) ShowAllPublishedRecipes(w http.ResponseWriter, r *http.Request) {

	recipes, err := h.App.DB.GetAllPublishedRecipes(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(recipes)
	if err != nil {
		return
	}
}

func (h *Handler) ShowFullRecipe(w http.ResponseWriter, r *http.Request) {

	//get id from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	//get recipe from the database
	recipe, err := h.App.DB.GetRecipeByID(r.Context(), id)
	log.Println("recipe error", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			log.Printf("ERROR: Could not get recipe by ID: %v", err)
			http.Error(w, "Failed to retrieve recipe", http.StatusInternalServerError)
		}
		return
	}

	//get ingredients from the database
	ingredients, err := h.App.DB.GetIngredientsByRecipeID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: Could not get ingredients: %v", err)
		http.Error(w, "Failed to retrieve ingredients in the recipe", http.StatusInternalServerError)
		// we don't return here, ingredients can be empty
	}
	// add ingredients to the recipe
	recipe.Ingredients = ingredients

	//get instructions from the database
	instructions, err := h.App.DB.GetInstructionsByRecipeID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: Could not get instructions: %v", err)
		http.Error(w, "Failed to retrieve instructions", http.StatusInternalServerError)
		// we don't return here, instructions can be empty
	}
	// add instructions to the recipe
	recipe.Instructions = instructions

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(recipe)
	if err != nil {
		return
	}
}

// ShowRecipesOfTheUser shows all recipes of logged-in user.
func (h *Handler) ShowRecipesOfTheUser(w http.ResponseWriter, r *http.Request) {

	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusInternalServerError)
		return
	}
	recipes, err := h.App.DB.GetAllRecipesFromUser(r.Context(), uid)
	if err != nil {
		http.Error(w, "Failed to retrieve recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(recipes)
	if err != nil {
		return
	}
	log.Println(recipes)
}

// GetAllIngredients fetch ingredients for create a new recipe page.
func (h *Handler) GetAllIngredients(w http.ResponseWriter, r *http.Request) {

	ingredients, err := h.App.DB.GetAllIngredients(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve ingredients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ingredients)
	if err != nil {
		return
	}

}

// SaveRecipe creates a new recipe and save it to the database.
func (h *Handler) SaveRecipe(w http.ResponseWriter, r *http.Request) {

	// save user id from session to the context
	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body into recipeForm
	var recipeForm model.RecipeForm
	err := json.NewDecoder(r.Body).Decode(&recipeForm)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate - check for mandatory fields
	if recipeForm.Title == "" || recipeForm.Portion <= 0 || len(recipeForm.Ingredients) < 1 || len(recipeForm.Instructions) < 1 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	//check that ingredients are not duplicated
	for i := 0; i < len(recipeForm.Ingredients); i++ {
		for j := i + 1; j < len(recipeForm.Ingredients); j++ {
			if recipeForm.Ingredients[i].IngredientId == recipeForm.Ingredients[j].IngredientId {
				http.Error(w, "Ingredients are duplicated", http.StatusBadRequest)
				return
			}
		}
	}

	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusInternalServerError)
		return
	}
	// save recipe in the DB
	recipe, err := h.App.DB.CreateRecipe(r.Context(), &recipeForm, uid)
	if err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(recipe)
	if err != nil {
		return
	}
}

// DeleteRecipe deletes a recipe from the database.
func (h *Handler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {

	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//get recipe id from the URL
	idStr := chi.URLParam(r, "id")
	recipeId, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusInternalServerError)
		return
	}
	err = h.App.DB.DeleteRecipe(r.Context(), recipeId, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) //204, no content
	log.Println("Recipe", recipeId, " deleted successfully")
}

type togglePublishPayload struct {
	Published *bool `json:"published"`
}

// UpdatePublishRecipe makes recipe visible for other users.
func (h *Handler) UpdatePublishRecipe(w http.ResponseWriter, r *http.Request) {

	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//get recipe id from the URL
	idStr := chi.URLParam(r, "id")
	recipeId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	// get wished recipe status from the request body
	var payload togglePublishPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Published == nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusInternalServerError)
		return
	}
	err = h.App.DB.PublishRecipe(r.Context(), recipeId, *payload.Published, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to update publish status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Publish status updated successfully")
	if err != nil {
		return
	}
}

// EditRecipe updates a recipe in the database.
func (h *Handler) EditRecipe(w http.ResponseWriter, r *http.Request) {

	userID := h.App.Session.Get(r.Context(), "userID")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//get recipe id from the URL
	idStr := chi.URLParam(r, "id")
	recipeId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var recipeForm model.RecipeForm
	// Parse request body into recipeForm
	err = json.NewDecoder(r.Body).Decode(&recipeForm)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//add recipeId to the recipeForm
	recipeForm.ID = recipeId

	// Validate - check for mandatory fields
	if recipeForm.Title == "" || recipeForm.Portion <= 0 || recipeForm.Ingredients == nil || len(recipeForm.Instructions) < 1 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	//check that ingredients are not duplicated
	for i := 0; i < len(recipeForm.Ingredients); i++ {
		for j := i + 1; j < len(recipeForm.Ingredients); j++ {
			if recipeForm.Ingredients[i].IngredientId == recipeForm.Ingredients[j].IngredientId {
				http.Error(w, "Ingredients are duplicated", http.StatusBadRequest)
				return
			}
		}
	}

	uid, ok := userID.(uint)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusInternalServerError)
		return
	}
	// save recipe in the DB
	err = h.App.DB.UpdateRecipe(r.Context(), &recipeForm, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Recipe not found or you do not have permission to edit it.", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Failed to update recipe: %v", err)
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Recipe updated successfully")
	if err != nil {
		return
	}

}

type adjustPortionPayload struct {
	NewServingSize int `json:"newServingSize"`
}

// AdjustServingSize change the amount of servings
func (h *Handler) AdjustServingSize(w http.ResponseWriter, r *http.Request) {

	//get recipe ID
	recipeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	// get servings size
	var payload adjustPortionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.NewServingSize <= 0 {
		http.Error(w, "Serving size needs to be greater than 0", http.StatusBadRequest)
		return
	}

	// get recipe from DB
	recipe, err := h.App.DB.GetRecipeByID(r.Context(), recipeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Printf("ERROR: Could not get recipe by ID: %v", err)
		http.Error(w, "Failed to retrieve recipe", http.StatusInternalServerError)
		return
	}

	ingredients, err := h.App.DB.GetIngredientsByRecipeID(r.Context(), recipeID)
	if err != nil {
		log.Printf("ERROR: Could not get ingredients: %v", err)
		http.Error(w, "Could not fetch ingredients", http.StatusInternalServerError)
		return
	}
	recipe.Ingredients = ingredients

	// call ingredients calculation
	newIngredients, err := service.CalculateIngredientsPerServing(recipe, payload.NewServingSize)
	if err != nil {
		log.Printf("ERROR: Could not calculate ingredients: %v", err)
		http.Error(w, "Could not calculate ingredients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newIngredients)
	if err != nil {
		return
	}
}
