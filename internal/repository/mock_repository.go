package repository

import (
	"context"

	"github.com/tjasha/Recipe_manager/internal/model"
)

// MockRepository is a mock repository for testing
type MockRepository struct {
	ExpectError               bool
	User                      *model.User
	Recipe                    model.Recipe
	Users                     []model.User
	SkipGoogleTokenValidation bool
}

// --- This is not used yet, just ready to start testing ---

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.ExpectError {
		return nil, context.DeadlineExceeded // returns error
	}
	// Returns the user
	return m.User, nil
}

func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
	if m.ExpectError {
		return 0, context.DeadlineExceeded
	}
	return 1, nil
}

func (m *MockRepository) GetAllPublishedRecipes(ctx context.Context) ([]model.Recipe, error) {
	return nil, nil
}
func (m *MockRepository) GetRecipeByID(ctx context.Context, id int64) (model.Recipe, error) {
	return m.Recipe, nil
}
func (m *MockRepository) GetIngredientsByRecipeID(ctx context.Context, id int64) ([]model.IngredientInRecipe, error) {
	return nil, nil
}
func (m *MockRepository) GetInstructionsByRecipeID(ctx context.Context, id int64) ([]model.RecipeInstruction, error) {
	return nil, nil
}
func (m *MockRepository) GetAllRecipesFromUser(ctx context.Context, userID uint) ([]model.Recipe, error) {
	return nil, nil
}
func (m *MockRepository) GetAllIngredients(ctx context.Context) ([]model.Ingredient, error) {
	return nil, nil
}
func (m *MockRepository) CreateRecipe(ctx context.Context, recipeData *model.RecipeForm, userID uint) (*model.Recipe, error) {
	return nil, nil
}
func (m *MockRepository) DeleteRecipe(ctx context.Context, recipeId int64, userId uint) error {
	return nil
}
func (m *MockRepository) PublishRecipe(ctx context.Context, recipeId int64, newStatus bool, userId uint) error {
	return nil
}
func (m *MockRepository) UpdateRecipe(ctx context.Context, recipeData *model.RecipeForm, userID uint) error {
	return nil
}

func (m *MockRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	if m.ExpectError {
		return nil, context.DeadlineExceeded
	}

	return m.Users, nil
}

func (m *MockRepository) DeleteUser(ctx context.Context, userId int) error {
	if m.ExpectError {
		return context.DeadlineExceeded
	}
	return nil
}

func (m *MockRepository) UpdateUserRole(ctx context.Context, userId, role int) error {
	if m.ExpectError {
		return context.DeadlineExceeded
	}
	return nil
}

func (m *MockRepository) UpdateUserState(ctx context.Context, userId int, state string) error {
	if m.ExpectError {
		return context.DeadlineExceeded
	}
	return nil
}
