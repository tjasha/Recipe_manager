package repository

import (
	"context"

	"github.com/tjasha/Recipe_manager/internal/model"
	"github.com/tjasha/Recipe_manager/internal/service"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllPublishedRecipes(ctx context.Context) ([]model.Recipe, error)
	GetRecipeByID(ctx context.Context, id int64) (model.Recipe, error)
	GetIngredientsByRecipeID(ctx context.Context, id int64) ([]model.IngredientInRecipe, error)
	GetInstructionsByRecipeID(ctx context.Context, id int64) ([]model.RecipeInstruction, error)
	GetAllRecipesFromUser(ctx context.Context, userID uint) ([]model.Recipe, error)
	GetAllIngredients(ctx context.Context) ([]model.Ingredient, error)
	CreateRecipe(ctx context.Context, recipeData *model.RecipeForm, userId uint) (*model.Recipe, error)
	DeleteRecipe(ctx context.Context, recipeId int64, userId uint) error
	PublishRecipe(ctx context.Context, recipeId int64, newStatus bool, userId uint) error
	UpdateRecipe(ctx context.Context, recipeData *model.RecipeForm, userId uint) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]model.User, error)
	DeleteUser(ctx context.Context, userId int) error
	UpdateUserRole(ctx context.Context, userId, role int) error
	UpdateUserState(ctx context.Context, userId int, state string) error
	GetFilteredAndSortedPublishedRecipes(ctx context.Context, filters service.ListRecipesParams) ([]model.Recipe, error)
}
