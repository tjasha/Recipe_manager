package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
		http.Error(w, "Cannot read token", http.StatusBadRequest)
		return
	}

	// Validate the token
	payload, err := idtoken.Validate(r.Context(), token.Credential, h.App.Config.GoogleOauthClientID)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// --- User logic ---
	email := payload.Claims["email"].(string)
	user, err := h.App.DB.GetUserByEmail(r.Context(), email)

	if err != nil {
		// If user does not exist, create a new one
		if errors.Is(err, pgx.ErrNoRows) {
			newUser := &model.User{
				UserName:    payload.Claims["name"].(string),
				Email:       email,
				GoogleID:    payload.Subject,
				AccessLevel: 1, // Default access level for new users
			}
			newID, createErr := h.App.DB.CreateUser(r.Context(), newUser)
			if createErr != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
			newUser.ID = uint(newID)
			user = newUser
		} else {
			// For any other database error
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Renew token to prevent session fixation attacks.
	err = h.App.Session.RenewToken(r.Context())
	if err != nil {
		http.Error(w, "Failed to renew session token", http.StatusInternalServerError)
		return
	}

	// Store user info in the session
	h.App.Session.Put(r.Context(), "userID", user.ID)
	h.App.Session.Put(r.Context(), "accessLevel", user.AccessLevel)

	// --- Respond to Frontend ---
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


//func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
//	id, err := h.service.CreateUser(context.Background(), "NewUser5", "NewUser5@test.com", "pass", 0)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//
//	fmt.Fprintf(w, "User created with ID:", id)
//
//	//fmt.Fprintf(w, "User created with ID:")
//
//}
//
//func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
//
//	//prevent session fixation attack
//	_ = m.App.Session.RenewToken(r.Context())
//
//	err := r.ParseForm()
//	if err != nil {
//		log.Println("Cannot parse form", err)
//	}
//
//	email := r.Form.Get("email")
//	password := r.Form.Get("password")
//
//	form := forms.New(r.PostForm)
//	form.Required("email", "password")
//	form.IsEmail("email")
//
//	// we want to send user back to the form
//	if !form.Valid() {
//		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
//			Form: form,
//		})
//		return
//	}
//
//	// try to authenticate the user
//	id, _, err := m.DB.Authenticate(email, password)
//	if err != nil {
//		log.Println("Authentication failed", err)
//		// if there is an error,i want to send user back to the log in form
//		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
//		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
//		return
//	}
//
//	// we authenticate a user by saving the ID that we got in the session
//	m.App.Session.Put(r.Context(), "user_id", id)
//	// i want to send user to home page after authentication
//	m.App.Session.Put(r.Context(), "flash", "Log in successful")
//	http.Redirect(w, r, "/", http.StatusSeeOther)
//
//}
//
//func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
//	_ = h.App.Session.Destroy(r.Context())
//	// renew session token
//	_ = h.App.Session.RenewToken(r.Context())
//
//	googleLogout()
//
//	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
//}
